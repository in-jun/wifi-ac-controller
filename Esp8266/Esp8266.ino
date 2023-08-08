#include <ESP8266WiFi.h>
#include <IRsend.h>
#include <ir_Samsung.h>
#include <IRremoteESP8266.h>
#include <ESP8266WebServer.h>
#include <ESP8266HTTPClient.h>

const int deviceIndex = 0; // 0~2
const char *ssid = "your-ssid";
const char *password = "your-password";

IRSamsungAc ac(4);
ESP8266WebServer server(80);

void sendDataToServer();
void handleOn();
void handleOff();

void setup()
{
    pinMode(LED_BUILTIN, OUTPUT);
    digitalWrite(LED_BUILTIN, LOW);

    WiFi.begin(ssid, password);
    WiFi.waitForConnectResult();

    sendDataToServer();

    ac.begin();
    ac.off();
    ac.setFan(kSamsungAcFanLow);
    ac.setMode(kSamsungAcCool);
    ac.setTemp(25);
    ac.setSwing(false);

    server.on("/on", handleOn);
    server.on("/off", handleOff);
    server.begin();

    digitalWrite(LED_BUILTIN, HIGH);
}

void sendDataToServer()
{
    WiFiClientSecure httpsClient;
    httpsClient.setInsecure();
    bool sendSuccess = false;

    while (!sendSuccess)
    {
        if (httpsClient.connect("aciotcontrol.onrender.com", 443))
        {
            String payload = "number=" + String(deviceIndex) + "&internal_ip=" + WiFi.localIP().toString();
            String request =
                String("POST ") + "/ip HTTP/1.1\r\n" +
                "Host: " + "aciotcontrol.onrender.com" + "\r\n" +
                "Content-Type: application/x-www-form-urlencoded\r\n" +
                "Content-Length: " + payload.length() + "\r\n" +
                "Connection: close\r\n\r\n" +
                payload + "\r\n";

            httpsClient.print(request);

            unsigned long timeout = millis();
            while (httpsClient.available() == 0 && millis() - timeout < 5000)
                delay(1);

            if (httpsClient.find("HTTP/1.1 200"))
                sendSuccess = true;

            httpsClient.stop();
            delay(1000);
        }
    }
}

void handleOn()
{
    ac.on();
    ac.send();
    server.send(200, "text/plain", "ON");
}

void handleOff()
{
    ac.off();
    ac.send();
    server.send(200, "text/plain", "OFF");
}

void loop()
{
    server.handleClient();
}
