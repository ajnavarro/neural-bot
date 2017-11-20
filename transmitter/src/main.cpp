#include <Arduino.h>
#include "RF24.h"

#define CE_PIN 7
#define CSN_PIN 8
#define MAX_SPEED = 400;

RF24 radio(CE_PIN, CSN_PIN);
const byte addresses[][6] = {"BOTCL", "FDBK"};

struct Request
{
    int leftSpeed;
    int rightSpeed;
};

struct Response
{
    int servoPosition;
    int distance;

    int leftSpeed;
    int rightSpeed;
};

void setup()
{
    Serial.begin(115200);
    radio.begin();

    radio.setAutoAck(true);
    radio.enableAckPayload();
    radio.enableDynamicPayloads();
    radio.setRetries(10, 5);
    radio.setChannel(125);
    radio.setDataRate(RF24_2MBPS);
    // Open writting pipe to ROBOT to send movement data
    radio.openWritingPipe(addresses[0]);
}

Request req;
Response res;
void loop()
{
    if (Serial.available() > 0)
    {
        String lString = Serial.readStringUntil(',');
        String rString = Serial.readStringUntil('|');

        req.leftSpeed = lString.toInt();
        req.rightSpeed = rString.toInt();

        bool ok = radio.write(&req, sizeof(req));
        if (!ok)
        {
            Serial.println("ERROR");
        }

        if (radio.isAckPayloadAvailable())
        {
            radio.read(&res, sizeof(Response));
            Serial.print(res.leftSpeed);
            Serial.print(",");
            Serial.print(res.rightSpeed);
            Serial.print(",");
            Serial.print(res.servoPosition);
            Serial.print(",");
            Serial.print(res.distance);
            Serial.print("|");
        }
    }
}