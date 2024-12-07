# elasticSearch
elasticSearch 구동하기 위한 가장 기본 적인 방법

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

---

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

---

3. `외부 접속 처리하기`

```
결국 우리는 외부에서 접속이 가능해야 한다. 이런 과정을 진행하기 위해서
yml 파일에 network.host: "_site_" 설정을 해주도록 하자

"_site_" 값은 ElasticSearch에서 다루는 변수 값으로
내부 IP 주소를 의미한다.

또한 추가로 "_site_" 값만 넣어주게 된다면, 외부에서만 접근이 가능하다.
그러기 떄문에 내부에서도 접근을 위해 "_local_" 설정도 기본으로 넣어주자

문제는 이렇게 설정을 하면, 외부에서도 접속이 가능한 상태로써 구동이 되기 떄문에
일반적으로 말하는 운영 환경의 설정으로 구동이 된다.

그래서 사실 이러게만 입력하고 실행을 하게 되면, 구동이 되지 않는다.
이후 여러가지 설정을 추가해 주어야 한다.
- 이떄 체크하는것이 부트스트랩 체크이며 대부분 이런 에러가 발생할 것이다.

1. process is too low
2. vm.max_map_count too low
3. discovery settings are unsuitable
```


한번 해결해 보자

> 1. process is too low
>
> 해당 에러는 사실 file descriptors에러이다.
ElasticSearch는 접근해야 하는 파일이 굉장히 많고, 그로인해서 하나의 프로세스가
여러개의 파일을 접근하여 처리하게 된다.
>
> 해당 프로세스가 접근 할 수 있는 파일의 제한량을 늘려줘야 하는 부분이다.
>
> https://www.elastic.co/guide/en/elasticsearch/reference/current/file-descriptors.html
> 
> 해당 링크를 참고하여 수정해주자.
> - ulimit는 일시적인 설정이기 떄문에 conf를 수정하는 방법으로 진행하자

> 2. vm.max_map_count too low
> 
> 해당 에러는 메모리가 너무 적기 떄문에 발생하는 에러이다
> 
> https://www.elastic.co/guide/en/elasticsearch/reference/current/vm-max-map-count.html
> 
> 해당 링크를 통해서 해결해주자

> 3. discovery settings ar unsuitable
> 
> 해당 설정은 단순하게 yaml 파일에서 discovery.seed_hosts
> 값에 host값만 넣어 주면 된다.
> 
> 추가로 해당 설정과 같이 가야 하는 정보는 cluster.initial_master_nodes 설정이다.
> 해당 값에도 node.name값이 동일하게 1:1 관계로 넣어주면 된다.
> 
> 예시는 다음과 같아.
> 
> > node.name: "node-1"
> >
> > discover.seed_hosts: ["test-1"]
> >
> > cluster.initial_master_nodes: ["node-1"]

해당 설정들이 마무리가 된다면,
- `sudo shutdown -r` 을 입력하여 인스턴스를 재실행 해주자


