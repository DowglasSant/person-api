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

// CreatePerson godoc
// @Summary      Create a new person
// @Description  Creates a new person in the system with the provided data
// @Tags         Persons
// @Accept       json
// @Produce      json
// @Param        person  body      contract.NewPersonDTO  true  "Person data to be created"
// @Success      201     {object}  contract.SuccessResponse
// @Failure      400     {object}  contract.ErrorResponse  "Invalid input data"
// @Failure      422     {object}  contract.ErrorResponse  "Business validation error"
// @Router       /persons [post]
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

// ListPersons godoc
// @Summary      List persons with pagination
// @Description  Returns a paginated list of persons with sorting options
// @Tags         Persons
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "Page number"              default(1)     minimum(1)
// @Param        page_size  query     int     false  "Items per page"           default(10)    minimum(1)  maximum(100)
// @Param        sort       query     string  false  "Field to sort by"         default(id)    Enums(id, name, cpf, email, created_at, updated_at)
// @Param        order      query     string  false  "Sort direction"           default(desc)  Enums(asc, desc)
// @Success      200        {object}  contract.PaginatedResponse{data=[]contract.PersonResponseDTO}
// @Failure      500        {object}  contract.ErrorResponse  "Internal server error"
// @Router       /persons [get]
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

// FindPersonByCPF godoc
// @Summary      Find person by CPF
// @Description  Returns person data based on the provided Brazilian CPF
// @Tags         Persons
// @Accept       json
// @Produce      json
// @Param        cpf  path      string  true  "Person's CPF (with or without formatting)"  example(111.444.777-35)
// @Success      200  {object}  contract.PersonResponseDTO
// @Failure      400  {object}  contract.ErrorResponse  "CPF not provided"
// @Failure      404  {object}  contract.ErrorResponse  "Person not found"
// @Failure      500  {object}  contract.ErrorResponse  "Internal server error"
// @Router       /persons/cpf/{cpf} [get]
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

// UpdatePerson godoc
// @Summary      Update a person
// @Description  Updates an existing person with the provided data
// @Tags         Persons
// @Accept       json
// @Produce      json
// @Param        id      path      int                     true  "Person ID"
// @Param        person  body      contract.UpdatePersonDTO  true  "Updated person data"
// @Success      200     {object}  contract.SuccessResponse
// @Failure      400     {object}  contract.ErrorResponse  "Invalid input data"
// @Failure      404     {object}  contract.ErrorResponse  "Person not found"
// @Failure      422     {object}  contract.ErrorResponse  "Business validation error"
// @Router       /persons/{id} [put]
func (h *PersonHandler) UpdatePerson(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("[ERROR] UpdatePerson - Invalid ID parameter: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Invalid person ID",
		})
		return
	}

	var dto contract.UpdatePersonDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Printf("[ERROR] UpdatePerson - Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Invalid request body: " + err.Error(),
		})
		return
	}

	log.Printf("[INFO] UpdatePerson - Updating person ID: %d", id)

	err = h.service.UpdatePerson(id, dto)
	if err != nil {
		if err.Error() == "person not found" {
			log.Printf("[WARN] UpdatePerson - Person not found with ID: %d", id)
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Person not found",
			})
			return
		}

		log.Printf("[ERROR] UpdatePerson - Failed to update person ID %d: %v", id, err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "validation_error",
			"message": err.Error(),
		})
		return
	}

	log.Printf("[SUCCESS] UpdatePerson - Person ID %d updated successfully", id)
	c.JSON(http.StatusOK, gin.H{
		"message": "Person updated successfully",
	})
}

// DeletePerson godoc
// @Summary      Delete a person
// @Description  Deletes a person from the system
// @Tags         Persons
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "Person ID"
// @Success      200 {object}  contract.SuccessResponse
// @Failure      400 {object}  contract.ErrorResponse  "Invalid ID"
// @Failure      404 {object}  contract.ErrorResponse  "Person not found"
// @Router       /persons/{id} [delete]
func (h *PersonHandler) DeletePerson(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("[ERROR] DeletePerson - Invalid ID parameter: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Invalid person ID",
		})
		return
	}

	log.Printf("[INFO] DeletePerson - Deleting person ID: %d", id)

	err = h.service.DeletePerson(id)
	if err != nil {
		if err.Error() == "person not found" {
			log.Printf("[WARN] DeletePerson - Person not found with ID: %d", id)
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Person not found",
			})
			return
		}

		log.Printf("[ERROR] DeletePerson - Failed to delete person ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to delete person: " + err.Error(),
		})
		return
	}

	log.Printf("[SUCCESS] DeletePerson - Person ID %d deleted successfully", id)
	c.JSON(http.StatusOK, gin.H{
		"message": "Person deleted successfully",
	})
}
