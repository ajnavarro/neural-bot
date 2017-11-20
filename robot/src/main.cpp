#include <Arduino.h>
#include "RF24.h"
#include <ZumoMotors.h>
#include <SimpleTimer.h>
#include "Servo.h"

#define CE_PIN 11
#define CSN_PIN 12
#define SERVO_PIN 5
#define SONAR_PIN 6

#define MAX_SPEED = 400;
#define DISTANCE 1

SimpleTimer timer;

RF24 radio(CE_PIN, CSN_PIN);
ZumoMotors motors;
Servo sonar;

const byte addresses[][6] = {"BOTCL", "FDBK"};

struct Response
{
    int leftSpeed;
    int rightSpeed;
};

struct Request
{
    int servoPosition;
    int distance;

    int leftSpeed;
    int rightSpeed;
};

Request req;
Response res;

void initRadio()
{
    radio.begin();

    radio.setAutoAck(true);
    radio.enableAckPayload();
    radio.setRetries(10, 5);
    radio.setChannel(125);
    radio.setDataRate(RF24_2MBPS);
    // Open reading pipe to read all the data sent to the robot
    radio.openReadingPipe(1, addresses[0]);
    radio.startListening();
}

void handleRadio()
{
    if (radio.available())
    {
        radio.read(&res, sizeof(Response));

        motors.setLeftSpeed(res.leftSpeed);
        motors.setRightSpeed(res.rightSpeed);

        req.leftSpeed = res.leftSpeed;
        req.rightSpeed = res.rightSpeed;

        radio.writeAckPayload(1, &req, sizeof(req));
    }
}

void initSonar()
{
    sonar.attach(SERVO_PIN, 20, 160);
    urm.begin(RX_SONAR_PIN, TX_SONAR_PIN, 9600);
}

const int sonarSteps[] = {0, 90, 180, 90};
const int arrayElements = 4;
int pos = 0;
int value;
void moveSonar()
{
    if (pos == arrayElements)
    {
        pos = 0;
    }

    req.servoPosition = sonarSteps[pos];
    sonar.write(req.servoPosition);
    resultType = urm.requestMeasurementOrTimeout(DISTANCE, value);
    if (resultType == DISTANCE)
    {
        req.distance = value;
    }
    else
    {
        req.distance = -1;
    }

    pos++;
}

void setup()
{
    initRadio();
    initSonar();

    timer.setInterval(500, moveSonar);
    timer.setInterval(10, handleRadio);
}

void loop()
{
    timer.run();
}