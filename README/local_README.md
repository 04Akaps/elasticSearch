Local에서는 파일을 다루기 귀찮으니, 그냥 Docker를 사용하자.

```
version: '3.6'
services:
  es01:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.12.2
    container_name: es01
    environment:
      - node.name=es01
      - bootstrap.memory_lock=true
      - "ES_JAVA_OPTS=-Xms1g -Xmx1g"
      - discovery.type=single-node
      - xpack.security.enabled=false
      - http.cors.enabled=true
      - http.cors.allow-origin="*"
      - http.cors.allow-methods=OPTIONS,HEAD,GET,POST,PUT,DELETE
      - http.cors.allow-headers=X-Requested-With,Content-Type,Content-Length,Authorization
      - http.cors.allow-credentials=true
    ulimits:
      memlock:
        soft: -1
        hard: -1
    volumes:
      - data01:/home/search/elastic_workspace
    ports:
      - 9200:9200
      - 9300:9300
    networks:
      - elastic
  kibana:
    image: docker.elastic.co/kibana/kibana:8.12.2
    container_name: kibana
    ports:
      - 5601:5601
    environment:
      - ELASTICSEARCH_HOSTS=["http://es01:9200"]
    depends_on:
      - es01
    networks:
      - elastic
volumes:
  data01:
    driver: local
networks:
  elastic:
    driver: bridge
```
