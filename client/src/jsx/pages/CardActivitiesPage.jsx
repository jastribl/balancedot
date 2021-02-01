import React, { useEffect, useState } from 'react'

import { getWithHandling, postForm } from '../../utils/api'
import ErrorRow from '../common/ErrorRow'
import Form from '../common/Form'
import Modal from '../common/Modal'
import Spinner from '../common/Spinner'
import CardActivitiesTable from '../tables/CardActivitiesTable'

const CardActivitiesPage = ({ match }) => {
    const cardUUID = match.params.cardUUID

    const [card, setCard] = useState(null)
    const [cardLoading, setCardLoading] = useState(false)
    const [modalVisible, setShowModal] = useState(false)
    const [errorMessage, setErrorMessage] = useState(null)

    const showModal = () => { setShowModal(true) }
    const hideModal = () => { setShowModal(false) }

    const refreshCard = () => getWithHandling(
        `/api/cards/${cardUUID}`,
        setCard,
        setErrorMessage,
        setCardLoading
    )

    const handleActivityUpload = (activityData) => {
        let formData = new FormData()
        for (let i = 0; i < activityData['files'].length; i++) {
            formData.append(`file${i}`, activityData['files'][i])
        }
        return postForm(`/api/cards/${cardUUID}/activities`, formData)
            .then(() => {
                hideModal()
                refreshCard()
            })
            .catch(e => {
                throw e.message
            })
    }

    useEffect(() => {
        refreshCard()
    }, [
        setCard,
        setErrorMessage,
        setCardLoading,
    ])

    return (
        <div>
            <Spinner visible={cardLoading} />
            <h1>Card Activities for {card ? (card.last_four + " (" + card.description + ")") : null}</h1>
            <input type='button' onClick={showModal} value='Upload Activities' style={{ marginBottom: 25 + 'px' }} />
            <ErrorRow message={errorMessage} />
            <CardActivitiesTable data={card?.activities ?? []} />
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
