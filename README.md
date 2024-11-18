# elasticSearch
elasticSearch 운영부터 query 작업


<h1> 운영 환경 </h1>

> AWS Amazon Linux


> *ElasticSearch Install*
>> wget https://www.elastic.co/kr/downloads/elasticsearch


>*Kibana Install*
>> wget https://artifacts.elastic.co/downloads/kibana/kibana-8.16.0-darwin-aarch64.tar.gz

`tar.gz`가 설치가 되었다면, `tar xfz` 명령어를 활용하여 풀어준다.
- `ElasticSearch`는 금방 풀리겠지만, `Kibana`는 파일이 좀 많아서 시간이 소요가 된다.



<h1> 기본 개념 </h1>

1. `ElasticSearch`가 실행이 되면, 기본적으로 `9200, 9300`을 사용한다.

```azure
이러한 이유는 기본적으로 여러 서버에 클러스터로 구동되는것을 목표로 하기 때문
   
클러스터로써 노드간 통신이 필요하기 때문에 
노드끼리 통신을 위한 포트와 Client 통신을 사용할 포트가 실행이 된다.
    
9300포트는 일반적으로 노드끼리의 통신을 위한 포트로 TCP를 활용
9200포트는 Client와 통신을 위한 포트로 HTTP 사용
   
해당 포트는 수정이 가능하다.
- config/elasticsearch.yml
```

2. `network.host`설정을 주의하자.

```azure
기본적인 Default값은 127.0.0.1이다.
또한 기본적으로 실행이 되면, 같은 호스트 에서만 접속이 가능하다.
    
하지만 ElasticSearch는 외부에서 접근이 가능해야 한다.
그래서 외부에서 접속을 할 수 있게, 실제 Network IP 주소를 넣어줘야한다.
    
하지만 이렇게 실제 Network IP주소가 들어가게 된다면, 그떄부터는 bootstrap을 체크하기 시작한다.
그래서 바로 실행이 안되는 이슈가 있다.
- 해당 부분은 나중에 cluster 설정하는 부분에서 추가적으로 알아본다.
    
    
network.bind_host : 내부망을 의미한다.
network.publish_host : 외부망을 의미한다.
```
