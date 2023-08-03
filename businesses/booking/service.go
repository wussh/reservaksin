package booking

import (
	"ca-reservaksin/businesses"
	"ca-reservaksin/businesses/session"
	"ca-reservaksin/helpers/nanoid"
	"strings"
)

type bookingsessionService struct {
	bookingRepository Repository
	sessionRepository session.Repository
}

func NewBookingSessionService(repoBooking Repository, repoSession session.Repository) Service {
	return &bookingsessionService{
		bookingRepository: repoBooking,
		sessionRepository: repoSession,
	}
}

func (service *bookingsessionService) BookingSession(dataBooking *Domain) (Domain, error) {
	dataBooking.Id, _ = nanoid.GenerateNanoId()

	getQueueNumber, err := service.bookingRepository.GetBySessionID(dataBooking.SessionId)
	if err != nil {
		return Domain{}, businesses.ErrInternalServer
	}
	dataBooking.NomorAntrian = len(getQueueNumber) + 1

	booking, err := service.bookingRepository.Create(dataBooking)
	if err != nil {
		return Domain{}, businesses.ErrInternalServer
	}

	getSessionByID, err := service.sessionRepository.GetByID(dataBooking.SessionId)
	if err != nil {
		return Domain{}, businesses.ErrInternalServer
	}

	getSessionByID.CapacityFulfilled += 1
	getSessionByID.Date = ""
	getSessionByID.StartSession = ""
	getSessionByID.EndSession = ""
	if _, err := service.sessionRepository.Update(dataBooking.SessionId, &getSessionByID); err != nil {
		return Domain{}, businesses.ErrInternalServer
	}

	return booking, nil
}

func (service *bookingsessionService) GetByCitizenID(citizenID string) ([]Domain, error) {
	dataBooking, err := service.bookingRepository.GetByCitizenID(citizenID)
	if err != nil {
		return []Domain{}, businesses.ErrInternalServer
	}

	return dataBooking, nil
}

func (service *bookingsessionService) GetBySessionID(sessionID string) ([]Domain, error) {
	dataBooking, err := service.bookingRepository.GetBySessionID(sessionID)
	if err != nil {
		return []Domain{}, businesses.ErrInternalServer
	}

	return dataBooking, nil
}

func (service *bookingsessionService) GetByNoKK(noKK string) ([]Domain, error) {
	dataBooking, err := service.bookingRepository.GetByNoKK(noKK)
	if err != nil {
		return []Domain{}, businesses.ErrInternalServer
	}

	return dataBooking, nil
}

func (service *bookingsessionService) UpdateStatusByID(id, status string) (Domain, error) {
	existed, err := service.bookingRepository.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return Domain{}, businesses.ErrIDNotFound
		}
		return Domain{}, businesses.ErrInternalServer
	}

	dataBooking, err := service.bookingRepository.UpdateStatusByID(existed.Id, status)
	if err != nil {
		return Domain{}, businesses.ErrInternalServer
	}

	if strings.ToLower(status) == "canceled" {
		getSessionByID, err := service.sessionRepository.GetByID(dataBooking.SessionId)
		if err != nil {
			return Domain{}, businesses.ErrInternalServer
		}

		getSessionByID.CapacityFulfilled -= 1
		getSessionByID.Date = ""
		getSessionByID.StartSession = ""
		getSessionByID.EndSession = ""
		if _, err := service.sessionRepository.Update(dataBooking.SessionId, &getSessionByID); err != nil {
			return Domain{}, businesses.ErrInternalServer
		}
	}

	return dataBooking, nil
}
