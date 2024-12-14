<h1>
 ElasticSearch API Server
</h1>

## ⚙ ElasticSearch 운영하기
> 기본적인 운영에 대한 방식은 README 폴더 참고


## ⚙ 기술 스택
> 해당 Project에서 사용된 디펜던시

`fx` : dependency injection를 활용하기 위해서 사용하였습니다.

`sonic` : web-server 특성상 성능적인 이점을 최대한 챙겨가기 위해서 직렬화 및 역직렬화에 적용을 하였습니다.

`fiber` : web-server의 api를 구성하며 빠른 성능을 보장하기 위해서 사용 되었습니다.

`redis` : 기본적인 Caching으로 사용이 되었으면 Cache Stamped 상황에 대해 방어하기 위한 PER 알고리즘이 들어가 있습니다.

`singleflight` : API 요청시에 불필요한 DB 처리에 대해 Process를 최적화 하기 위해 사용 되었습니다.

## ⚙ 사용된 캐싱 전략 및 Process 
> `/common/strategy` 참고

1. *PER algorithm*

Cache Hit을 주기적으로 유도하기 위해 사용이 되었고, Cache stamped 방어를 위해 사용 되었습니다.

2. *Singleflight*

대용량 트래픽 상황에 대해 DB에 대한 접근은 최소화 하기 위해서 API Process를 키값을 기준으로
단 하나의 요청만 허요하고 처리하는 개념으로써 사용 되었씁니다.