import React, { useEffect, useState } from 'react'

import { get, postForm } from '../../utils/api'
import Form from '../common/Form'
import Modal from '../common/Modal'
import CardActivitiesTable from '../tables/CardActivitiesTable'

const CardActivitiesPage = ({ match }) => {
    const cardUUID = match.params.cardUUID

    const [card, setCard] = useState(null)
    const [cardActivities, setCardActivities] = useState([])
    const [modalVisible, setShowModal] = useState(false)

    const showModal = () => { setShowModal(true) }
    const hideModal = () => { setShowModal(false) }

    const refreshCard = () => {
        get(`/api/cards/${cardUUID}`)
            .then(cardResponse => setCard(cardResponse))
    }

    const refreshCardActivities = () => {
        get(`/api/cards/${cardUUID}/activities`)
            .then(cardActivities => setCardActivities(cardActivities))
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
            <h1>Card Activities for {card ? (card.last_four + " (" + card.description + ")") : null}</h1>
            <input type='button' onClick={showModal} value='Upload Activities' style={{ marginBottom: 25 + 'px' }} />
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
