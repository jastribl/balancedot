import React, { useEffect, useState } from 'react'

import { get, postForm } from '../../utils/api'
import ErrorRow from '../common/ErrorRow'
import Form from '../common/Form'
import Modal from '../common/Modal'
import Spinner from '../common/Spinner'
import CardActivitiesTable from '../tables/CardActivitiesTable'

const CardActivitiesPage = ({ match }) => {
    const cardUUID = match.params.cardUUID

    const [card, setCard] = useState(null)
    const [cardActivities, setCardActivities] = useState([])
    const [cardsLoading, setCardsLoading] = useState(false)
    const [cardActivitiesLoading, setCardActivitiesLoading] = useState(false)
    const [modalVisible, setShowModal] = useState(false)
    const [errorMessage, setErrorMessage] = useState(null)

    const showModal = () => { setShowModal(true) }
    const hideModal = () => { setShowModal(false) }

    const refreshCard = () => {
        setCardsLoading(true)
        get(`/api/cards/${cardUUID}`)
            .then(cardResponse => setCard(cardResponse))
            .catch(e => setErrorMessage(e.message))
            .finally(() => setCardsLoading(false))
    }

    const refreshCardActivities = () => {
        setCardActivitiesLoading(true)
        get(`/api/cards/${cardUUID}/activities`)
            .then(cardActivities => setCardActivities(cardActivities))
            .catch(e => setErrorMessage(e.message))
            .finally(() => setCardActivitiesLoading(false))
    }

    const handleActivityUpload = (activityData) => {
        let formData = new FormData()
        for (let i = 0; i < activityData['files'].length; i++) {
            formData.append(`file${i}`, activityData['files'][i])
        }
        return postForm(`/api/cards/${cardUUID}/activities`, formData)
            .then(() => {
                hideModal()
                refreshCardActivities()
            })
            .catch(e => {
                throw e.message
            })
    }

    useEffect(() => {
        refreshCard()
        refreshCardActivities()
    }, [setCardActivities])

    return (
        <div>
            <Spinner visible={cardsLoading || cardActivitiesLoading} />
            <h1>Card Activities for {card ? (card.last_four + " (" + card.description + ")") : null}</h1>
            <input type='button' onClick={showModal} value='Upload Activities' style={{ marginBottom: 25 + 'px' }} />
            <ErrorRow message={errorMessage} />
            <CardActivitiesTable data={cardActivities} />
            <Modal headerText='Activity Upload' visible={modalVisible} handleClose={hideModal}>
                <Form
                    onSubmit={handleActivityUpload}
                    fieldInfos={{
                        files: {
                            inputType: 'file',
                            multiple: true,
                        },
                    }}
                />
            </Modal>
        </div>
    )
}

export default CardActivitiesPage
