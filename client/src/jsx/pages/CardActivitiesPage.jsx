import Moment from 'moment'
import React, { useEffect, useState } from 'react'

import { get, postForm } from '../../utils/api'
import { formatAsMoney } from '../../utils/format'
import Form from '../common/Form'
import Modal from '../common/Modal'
import Table from '../common/Table'

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
        let formData = new FormData();
        formData.append('file', activityData['file'])
        return postForm(`/api/cards/${cardUUID}/activity`, formData)
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
            <h1>Card Activities for {card?.last_four}</h1>
            <input type='button' onClick={showModal} value='Upload Activities' style={{ marginBottom: 25 + 'px' }} />
            <div>
                <Table rowKey='uuid' columns={{
                    'uuid': 'Activity UUID',
                    'transaction_date': 'Transaction Date',
                    'post_date': 'Post Date',
                    'description': 'Description',
                    'category': 'Category',
                    'type': 'Type',
                    'amount': 'Amount',
                }} rows={cardActivities} customRenders={{
                    'transaction_date': (data) =>
                        Moment(data['transaction_date']).format('YYYY-MM-DD'),
                    'post_date': (data) =>
                        Moment(data['post_date']).format('YYYY-MM-DD'),
                    'amount': (data) => formatAsMoney(data['amount']),
                }} />
            </div>
            <Modal headerText='Activity Upload' visible={modalVisible} handleClose={hideModal}>
                <Form
                    onSubmit={handleActivityUpload}
                    fieldInfos={{
                        file: {
                            fieldLabel: 'File',
                            fieldName: 'file',
                            placeholder: 'File...',
                            inputType: 'file',
                        },
                    }}
                />
            </Modal>
        </div>
    )
}

export default CardActivitiesPage
