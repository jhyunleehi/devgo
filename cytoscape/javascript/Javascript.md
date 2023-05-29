# Javascript

자바스크립트는 브라우져에서 실행될때 마다 해서되는 인터프리터형의 스크립트 프로그래밍 언어로 동적인페이지를만들기 위해서 사용된다. 웹 뿐만 아니라 모바일 웹, 서버,  IOT 등 다양한 분양에서 사용이 가능하다.

## javascript 사용 방법
### HTML 문서
* 코드에 직접 입력하는 방법
```
<script type ="text/javascript"> ...Code... </script>
```
* 외부 script 연동
```
<script src="./js/example.js" type="text/javascript"></script>
```
*  주석 처리
```
//
/* ~ */
```

### 변수
#### 변수 타입 
java, c, c++ 과는 다르게 명시적인 타입이 없는 느슨한 데이터  type 어떤 자료형(문자열, 숮나, 객체, 함수) 값도 담을 수 있고 한 변수를 다 type의 값으로  할당 할 수 있다. 
#### 변수 위치
변수의 위치에 따라서 scope가 달라지기 때문에 주의
#### global local 변수
local 변수는 var로 시작한다. var로 시작하지 않으면 전역 변수로 인식한다.
local 변수는 함수 안에서 선언하고 global 변수와 local 변수가 충돌이 발생하는 경우 local 변수가 우선한다. 

#### 클로져 Closure
자바스크립트의 모든 함수는 클로저이다. 내부 함수가 외부 함수의 context에 접근 할 수 있고, 외부 함수의 실행이 끝나서 외부 함수가 소멸된 이후에도 내부 함수가 외부함수의 변수에 접근할 수 있는 메커니즘이다.
```
var sequencer = function() {
    var s = 0;
    return function() {
        return ++s;
    }
}

var seq = sequencer();

seq();  // 1
seq();  // 2
seq();  // 3
```
```
var items = document.getElementsByTagName(‘li’);
for (var i = 0; i < items.length; i++) {
    (function() {   // 새로운 스코프 선언
        var idx = i;    // 클로저가 접근할 수 있는 변수 선언
        items[i].onclick = function(e) {
        console.log(‘My Sequence is ‘ + (i + 1);    // 자신의 순번 출력
        }
    })();
}
```
## Node.js
웹개발은 크게 둘로 나뉩니다. 프론트엔드 와 백엔트 기술입니다.  Node.js 등장으로 백엔드까지 javascript가 사용할 수 있게 되었다.
### 특징
* javascript 엔진에서 동작하는 이벤트 처리 I/O 프레임워크
* 웹서버와 같이 확장성 있는 네트워크 프로그램 작성을 위해서 고안
* 서버 측 실행

#### Node.js 예시
hello world http 서버

```
 var http = require('http');
 http.createServer(function ( request, response) {
    response.writeHead(200, {'Content-Type' : 'text/plain'});
    response.end('Hello World\n');
 }).listen(8000);
 console.log('Server running at http://localhost:8000/');
 ```

 ```
 const http = require('http');
 const hostname = '127.0.0.1';
 const port = 3000;
 const server = http.createServer((req,res) => {
    res.statusCode = 200
    res.setHeader('Content-type', 'text/plain')
    res.end('Hello World\n')
 })
 server.listen(port, hostname, () =>{
    console.log('server running at http://${hostname}:${port}/')
 })
 ```
 ```
 var net = require('net');
 net.createServer(function(stream){
    stream.write('hello\r\n');
    stream.on('end', function(){
        streamm.end('goodbye\r\n');
        });
    stream.pipe(stream);
 }).listen(7000);
 ```
### Node.js npm 설치 : windows

<https://hello-bryan.tistory.com/95>

### npm

npm은 Node Packaged Manager의 약자입니다. 먼저 Node는 Node.js를 의미하는 것 같습니다. Packaged라는 것은 package로 만들어진 것들을 의미하는 것 같습니다. package는 모듈이라고도 불리는데 패키지나 모듈은 프로그램보다는 조금 작은 단위의 기능들을 의미합니다. 그리고 Manager는 잘 아시는 것처럼 관리자를 의미합니다.  이걸 합쳐보면 npm이라는 것은 Node.js로 만들어진 pakage(module)을 관리해주는 툴이라는 것이 됩니다. 

이름처럼 npm은 Node.js로 만들어진 모듈을 웹에서 받아서 설치하고 관리해주는 프로그램입니다. 개발자는 단 몇 줄의 명령어로 기존에 공개된 모듈들을 설치하고 활용할 수 있습니다. 프로그램보다 조금 작은 단위인 이 모듈들을 필요에 따라서 이런 저런 모양으로 쌓아서 활용을 할 수 있다고 하는데 필요한 기능을 적절하게 활용할 수 있다면 개발자 입장에서는 참 좋은 일이죠(Java랑 비교를 하자면 메이븐과 비슷한 역할을 하는 것 같습니다).

거기서 그치는 것이 아니라 이 모듈들을 활용했다면 이후에 그 모듈을 만든 개발자가 업데이트를 하거나 할 경우 체크를 해서 알려주는 듯합니다. 버전관리도 쉬워진다는 의미입니다. 

#### npm 사용
webpack package 설치
```
npm install --g webpack
npm init
```
package.json 파일이 자동으로 생성
```
{
  "name": "application-study",
  "version": "1.0.0",
  "description": "",
  "main": "index.js",
  "dependencies": {
    "http-server": "0.9.0",
    "rimraf": "2.6.1",
    "webpack": "2.2.1",
    "worker-loader": "0.8.0"
  },
  "scripts": {
    "prebuild": "rimraf dist",
    "build": "webpack --config webpack/webpack.config.js",
    "http-server": "http-server -c-1"
  },
  "repository": {
    "type": "git",
    "url": "git+https://github.com/aaa/bbb.git"
  },
  "author": "",
  "license": "ISC",
  "bugs": {
    "url": "https://github.com/aaa/bbb/issues"
  },
  "homepage": "https://github.com/aaa/bbb#readme"
}
```
package.json의 script 구문 실행 방법

webpack --config webpack/webpack.config.js 실행 방법
```
npm run build
npm run http-server

```

#### npm 명령
```
$ npm

Usage: npm <command>

where <command> is one of:
    access, adduser, audit, bin, bugs, c, cache, ci, cit,       
    clean-install, clean-install-test, completion, config,      
    create, ddp, dedupe, deprecate, dist-tag, docs, doctor,     
    edit, explore, fund, get, help, help-search, hook, i, init, 
    install, install-ci-test, install-test, it, link, list, ln, 
    login, logout, ls, org, outdated, owner, pack, ping, prefix,
    profile, prune, publish, rb, rebuild, repo, restart, root,  
    run, run-script, s, se, search, set, shrinkwrap, star,      
    stars, start, stop, t, team, test, token, tst, un,
    uninstall, unpublish, unstar, up, update, v, version, view,
    whoami

npm <command> -h  quick help on <command>
npm -l            display full usage info
npm help <term>   search for help on <term>
npm help npm      involved overview



jhyunlee@SUPER-COM MINGW64 /c/Gocode/src/github.com/jhyunleehi/go (master)
$ npm install bower
npm WARN deprecated bower@1.8.8: We don't recommend using Bower for new projects. Please consider Yarn and Webpack or Parcel. You can read how to migrate legacy project here: https://bower.io/blog/2017/how-to-migrate-away-from-bower/
npm WARN saveError ENOENT: no such file or directory, open 'C:\Gocode\src\github.com\jhyunleehi\go\package.json'
npm notice created a lockfile as package-lock.json. You should commit this file.
npm WARN enoent ENOENT: no such file or directory, open 'C:\Gocode\src\github.com\jhyunleehi\go\package.json'
npm WARN go No description
npm WARN go No repository field.
npm WARN go No README data
npm WARN go No license field.

+ bower@1.8.8
added 1 package from 1 contributor and audited 1 package in 3.97s
found 0 vulnerabilities

```



## AngularJS
 javascript로 구현한 Client측 MVC/MVVM 프레임워크 

 ### MVC
* Model : 보통 JSON으로 표현되는 애플리케이션의 특정한 데이터 구조를 말한다.
* View : HTML 혹은 렌더링된 결과를 말한다. MVC 프레임워크를 사용한다면 뷰를 갱신할 모델 데이터를 내려받은 뒤 HTML에서 해당 데이터를 보여줄 것이다.
* Controller 서버 에서 직접 뷰 로 접근하는 일종의 중간 통로로서 필요할 때마다 서버와 클라이언트 통신으로 데이터를 변경한다.