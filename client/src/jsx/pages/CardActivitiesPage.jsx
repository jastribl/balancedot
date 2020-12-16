import React, { useEffect, useState } from 'react'
import Moment from 'moment'

import { post, get } from '../../utils/api'

import Table from "../common/Table"
import Modal from "../common/Modal"
import Form from "../common/Form"

const CardActivitiesPage = ({ match }) => {
    const cardUUID = match.params.cardUUID

    const [card, setCard] = useState(null)
    const [cardActivities, setCardActivities] = useState([])
    // const [modalVisible, setShowModal] = useState(false)

    // const showModal = () => { setShowModal(true) }
    // const hideModal = () => { setShowModal(false) }

    const refreshCard = () => {
        get(`/api/cards/${cardUUID}`)
            .then((card) => setCard(card))
    }

    const refreshCardActivities = () => {
        get(`/api/cards/${cardUUID}/activities`)
            .then((cardActivities) => setCardActivities(cardActivities))
    }

    // const handleNewCardSubmit = (newCardData) => {
    //     return post('/api/card', newCardData)
    //         .then(() => {
    //             hideModal()
    //             refreshCardActivities()
    //         })
    //         .catch(e => {
    //             throw e.message
    //         })
    // }

    useEffect(() => {
        refreshCard()
        refreshCardActivities()
    }, [setCardActivities])

    return (
        <div>
            <h1>Card Activities for {card?.last_four}</h1>
            {/* <input type="button" onClick={showModal} value="New Card" style={{ marginBottom: 25 + 'px' }} /> */}
            <div>
                <Table rowKey="uuid" columns={{
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
                    'amount': (data) => (data['amount'] < 0 ? '-' : '') + '$' + Math.abs(data['amount'])
                }} />
            </div>
            {/* <Modal headerText="New Card" visible={modalVisible} handleClose={hideModal}>
                <Form
                    onSubmit={handleNewCardSubmit}
                    fieldInfos={{
                        last_four: {
                            fieldLabel: "Last Four",
                            fieldName: "last_four",
                            placeholder: "Last Four...",
                            inputType: "text",
                            validate: (fieldLabel, fieldValue) => {
                                if (!/^[0-9][0-9][0-9][0-9]$/.test(fieldValue)) {
                                    return `${fieldLabel} must follow '####' format`
                                }
                                return null
                            }
                        },
                        description: {
                            fieldLabel: "Description",
                            fieldName: "description",
                            placeholder: "Description...",
                            inputType: "text",
                            validate: (fieldLabel, fieldValue) => {
                                if (!/...../.test(fieldValue)) {
                                    return `${fieldLabel} must be at least 4 characters long`
                                }
                                return null
                            }
                        }
                    }}
                />
            </Modal> */}
        </div>
    )
}

export default CardActivitiesPage
