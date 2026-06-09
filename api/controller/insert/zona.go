package insert

import (
	"errors"
	"log"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/rafinhacuri/SanchezDNS/api/zonas"
)

type Soa struct {
	StartOfAuthority string `json:"startOfAuthority"`
	Email            string `json:"email"`
	Refresh          int    `json:"refresh"`
	Retry            int    `json:"retry"`
	Expire           int    `json:"expire"`
	NegativeCacheTtl int    `json:"negativeCacheTtl"`
}

func (s *Soa) Validate() error {
	if s.StartOfAuthority == "" {
		return errors.New("start of authority é obrigatório")
	}

	if s.Email == "" {
		return errors.New("email é obrigatório")
	}

	if s.Refresh <= 0 {
		return errors.New("refresh deve ser um inteiro positivo")
	}

	if s.Retry <= 0 {
		return errors.New("retry deve ser um inteiro positivo")
	}

	if s.Expire <= 0 {
		return errors.New("expire deve ser um inteiro positivo")
	}

	if s.NegativeCacheTtl <= 0 {
		return errors.New("negative cache ttl deve ser um inteiro positivo")
	}

	return nil
}

type CreateZoneRequest struct {
	Domain string `binding:"required" json:"domain"`
	Soa    Soa    `binding:"required" json:"soa"`
	Type   string `binding:"required" json:"type"`
}

func (req *CreateZoneRequest) Validate() error {
	if req.Domain == "" {
		return errors.New("domain é obrigatório")
	}

	if !strings.HasSuffix(req.Domain, ".") {
		req.Domain += "."
	}

	req.Type = strings.ToLower(strings.TrimSpace(req.Type))
	switch req.Type {
	case "", "normal", "forward":
		req.Type = "normal"
	case "reverse":
		if !strings.HasSuffix(req.Domain, ".in-addr.arpa") {
			return errors.New("reverse zone deve terminar com .in-addr.arpa")
		}
	case "reverse-ipv6":
		if !strings.HasSuffix(req.Domain, ".ip6.arpa") {
			return errors.New("reverse-ipv6 zone deve terminar com .ip6.arpa")
		}
	default:
		return errors.New("type deve ser um dos seguintes: normal, reverse, reverse-ipv6")
	}

	err := req.Soa.Validate()
	if err != nil {
		return errors.New("soa: " + err.Error())
	}

	return nil
}

func Zone(c *gin.Context) {
	var req CreateZoneRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Requisição inválida"})

		return
	}

	err = req.Validate()
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "erro de validação"})

		return
	}

	ctx := c.Request.Context()

	email := c.GetString("email")

	_, err = zonas.CreateZone(
		ctx,
		req.Domain,
		req.Type,
		req.Soa.StartOfAuthority,
		req.Soa.Email,
		req.Soa.Refresh,
		req.Soa.Retry,
		req.Soa.Expire,
		req.Soa.NegativeCacheTtl,
		email)
	if err != nil {
		log.Println(err)

		c.AbortWithStatusJSON(502, gin.H{"message": "falha ao criar zona"})

		return
	}

	c.JSON(201, gin.H{"message": "zona criada com sucesso"})
}
