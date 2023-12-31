# 에어컨 원격 제어 프로젝트

이 프로젝트는 ESP8266 마이크로컨트롤러 보드를 활용하여 에어컨을 웹을 통해 원격으로 제어하는 흥미로운 프로젝트입니다.

## 필요한 라이브러리

프로젝트를 구현하기 위해서는 아래의 라이브러리들이 필요합니다. Arduino IDE의 라이브러리 관리자에서 각 라이브러리 이름을 검색하여 쉽게 설치할 수 있습니다.

1. ESP8266WiFi
2. IRsend
3. ir_Samsung
4. IRremoteESP8266
5. ESP8266WebServer
6. ESP8266HTTPClient

## 기능 및 특징

이 프로젝트는 몇 가지 멋진 기능과 특징을 가지고 있습니다:

-   웹 브라우저를 통해 에어컨을 간편하게 켜고 끌 수 있습니다.
-   각각의 디바이스에 고유한 deviceIndex를 할당하여 다중 제어가 가능합니다.
-   다수의 디바이스를 함께 제어하여 사용자의 경험을 향상시킵니다.
-   Golang 웹 서버를 통해 클라이언트의 내부 IP 주소 매핑 및 리디렉션을 간편하게 관리합니다.

## 설치 및 사용 방법

프로젝트를 시작하는 방법을 살펴보겠습니다:

1. 저장소를 클론합니다: `git clone https://github.com/in-jun/wifi-ac-controller`
2. Arduino IDE를 열고 `Esp8266.ino` 파일을 엽니다.
3. 필요한 라이브러리를 설치합니다. 위에서 언급한 라이브러리들은 Arduino IDE의 라이브러리 관리자에서 손쉽게 설치할 수 있습니다.
4. Wi-Fi 정보를 입력합니다:

    ```cpp
    const char *ssid = "your-ssid";
    const char *password = "your-password";
    ```

5. 디바이스에 맞게 설정을 변경합니다:

    ```cpp
    const int deviceIndex = 0; // 디바이스의 인덱스
    ```

6. 업로드 버튼을 눌러 ESP8266 보드에 프로그램을 업로드합니다.
7. 웹 브라우저에서 `https://aciotcontrol.onrender.com/${deviceIndex}/on`과 같은 URL을 입력하여 에어컨을 쉽게 제어하세요.

## ESP8266 모듈과 IR LED 연결

이 프로젝트에서는 에어컨을 제어하기 위해 IR LED를 활용하며, 이를 위해 ESP8266 모듈과 IR LED를 연결해야 합니다. 연결 방법을 간략히 알아보겠습니다:

-   IR LED 숏 핀 (빨간색): ESP8266의 GPIO 핀 (D4)
-   IR LED 롱 핀 (흰색): GND (지상, 그라운드) 핀

IR LED의 숏 핀을 ESP8266의 GPIO 핀에 연결하고, 롱 핀을 GND에 연결합니다. 이렇게 연결하면 ESP8266 모듈로부터 IR LED를 제어할 수 있습니다.

## 내부 LED 활용

ESP8266 모듈의 내부 LED를 활용하여 초기화 및 에러 상태를 나타낼 수 있습니다. `setup()` 함수 내에서 내부 LED를 제어하는 코드가 포함되어 있습니다.

## Golang 웹 서버

### 요청 처리 흐름

클라이언트의 요청을 처리하는 과정은 다음과 같습니다:

1. 클라이언트가 특정 URL을 요청합니다.
2. 웹 서버는 받은 URL에서 디바이스 인덱스와 경로를 파악합니다.
3. 서버는 클라이언트의 공용 IP 주소를 식별합니다.
4. 내부 IP 매핑 정보를 참고하여 해당 공용 IP 주소에 대한 내부 IP 주소를 확인합니다.
5. 확인된 내부 IP 주소로 리디렉션 응답을 생성하고, 클라이언트에게 보냅니다.
6. 클라이언트는 서버로부터 받은 리디렉션 응답을 처리하여 해당 내부 IP 주소로 접속합니다.

이렇게 함으로써 클라이언트는 ESP8266 모듈의 내부 IP 주소를 알 필요 없이, 서버를 통해 쉽게 ESP8266 모듈에 접속할 수 있습니다. 이는 내부 IP 주소 노출 없이도 디바이스에 편리하게 접속할 수 있는 장점으로 사용자 경험을 향상시킵니다.
