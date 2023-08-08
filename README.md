# 에어컨 원격 제어 프로젝트

이 프로젝트는 ESP8266 마이크로컨트롤러 보드를 활용하여 웹을 통해 에어컨을 원격으로 제어하는 기능을 구현한 프로젝트입니다.

## 필요한 라이브러리

이 프로젝트를 위해 다음 라이브러리들이 필요합니다. 아래의 각 라이브러리 이름을 Arduino IDE의 라이브러리 관리자에서 검색하여 설치하세요.

1. ESP8266WiFi
2. IRsend
3. ir_Samsung
4. IRremoteESP8266
5. ESP8266WebServer
6. ESP8266HTTPClient

## 기능 및 특징

-   웹 브라우저를 통해 에어컨을 켜고 끌 수 있습니다.
-   각각의 디바이스에 deviceIndex 를 할당하여 다중 제어 가능합니다.
-   여러 디바이스의 사용을 허용하여 사용자의 편의성을 증대시킵니다.

## 설치 및 사용 방법

1. 저장소를 클론합니다: `git clone https://github.com/in-jun/wifi-ac-controller`
2. Arduino IDE를 열고 `Esp8266.ino` 파일을 엽니다.
3. 필요한 라이브러리를 설치합니다: 위에서 언급한 라이브러리들.
4. Wi-Fi 정보를 입력합니다:

    ```cpp
    const char *ssid = "your-ssid";
    const char *password = "your-password";
    ```

5. 디바이스에 맞게 설정을 변경합니다:
    ```cpp
    const int deviceIndex = 0; // 디바이스의 인덱스
    ```
6. ESP8266 보드에 프로그램을 업로드합니다.
7. 웹 브라우저를 열고 `https://aciotcontrol.onrender.com/${deviceIndex}/on`과 같은 URL을 입력하여 에어컨을 제어합니다.
