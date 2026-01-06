package handler

import (
	"log"
	"math"
	"net/http"
	"strconv"

	contract "pessoas-api/internal/contract/person"
	"pessoas-api/internal/domain/person/ports"

	"github.com/gin-gonic/gin"
)

type PersonHandler struct {
	service ports.PersonService
}

func NewPersonHandler(service ports.PersonService) *PersonHandler {
	return &PersonHandler{
		service: service,
	}
}

func (h *PersonHandler) CreatePerson(c *gin.Context) {
	var dto contract.NewPersonDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Printf("[ERROR] CreatePerson - Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Invalid request body: " + err.Error(),
		})
		return
	}

	id, err := h.service.CreatePerson(dto)
	if err != nil {
		log.Printf("[ERROR] CreatePerson - Validation error for CPF %s: %v", dto.CPF, err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "validation_error",
			"message": err.Error(),
		})
		return
	}

	log.Printf("[SUCCESS] CreatePerson - Person created with ID: %d, Name: %s", id, dto.Name)
	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"message": "Person created successfully",
	})
}

func (h *PersonHandler) ListPersons(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	sort := c.DefaultQuery("sort", "id")
	order := c.DefaultQuery("order", "desc")

	log.Printf("[INFO] ListPersons - Fetching page: %d, pageSize: %d, sort: %s, order: %s", page, pageSize, sort, order)

	persons, total, err := h.service.ListPersons(page, pageSize, sort, order)
	if err != nil {
		log.Printf("[ERROR] ListPersons - Failed to retrieve persons: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to retrieve persons: " + err.Error(),
		})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	log.Printf("[SUCCESS] ListPersons - Retrieved %d persons (total: %d, pages: %d)", len(persons), total, totalPages)

	response := contract.PaginatedResponse{
		Data:       persons,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, response)
}

func (h *PersonHandler) FindPersonByCPF(c *gin.Context) {
	cpf := c.Param("cpf")

	if cpf == "" {
		log.Printf("[ERROR] FindPersonByCPF - CPF parameter is empty")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "CPF parameter is required",
		})
		return
	}

	log.Printf("[INFO] FindPersonByCPF - Searching for CPF: %s", cpf)

	person, err := h.service.FindPersonByCPF(cpf)
	if err != nil {
		log.Printf("[ERROR] FindPersonByCPF - Failed to find person with CPF %s: %v", cpf, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to find person: " + err.Error(),
		})
		return
	}

	if person == nil {
		log.Printf("[WARN] FindPersonByCPF - Person not found with CPF: %s", cpf)
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "not_found",
			"message": "Person not found with the provided CPF",
		})
		return
	}

	log.Printf("[SUCCESS] FindPersonByCPF - Found person with ID: %d, Name: %s", person.ID, person.Name)
	c.JSON(http.StatusOK, person)
}
