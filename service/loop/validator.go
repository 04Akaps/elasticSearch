package loop

import (
	"github.com/04Akaps/elasticSearch.git/common/validator"
	"github.com/04Akaps/elasticSearch.git/config"
	"github.com/04Akaps/elasticSearch.git/repository/elasticSearch"
	_validator "github.com/go-playground/validator"
	"log"
	"time"
)

type ValidatorLoop struct {
	cfg           config.Config
	ElasticSearch elasticSearch.ElasticSearch
}

func RunValidatorLoop(
	cfg config.Config,
	elasticSearch elasticSearch.ElasticSearch,
) {
	l := ValidatorLoop{cfg, elasticSearch}

	go l.RegisterValidation()

	elasticSearch.Indexes()
}

func (v ValidatorLoop) RegisterValidation() {

	_t := time.NewTicker(10 * time.Second)
	defer _t.Stop()

	for {
		v.indexValidation()

		<-_t.C
	}
}

func (v ValidatorLoop) indexValidation() {
	res, err := v.ElasticSearch.Indexes()

	if err != nil {
		log.Println("Failed to get indexes", "err", err)
		return
	}

	// -> elastic에서 설정한 index 요청한 값만 request 레벨에서 검증하기 위한 벨리데이션
	validator.RegisterValidation("index", func(f1 _validator.FieldLevel) bool {
		value := f1.Field().String()
		return res[value]
	})
}
