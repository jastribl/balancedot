import React, { useEffect, useState } from 'react'

import { post, get } from '../../utils/api'

import Table from "../common/Table"
import Modal from "../common/Modal"
import Form from "../common/Form"

const CardsPage = () => {
    const [cards, setCards] = useState(null)
    const [modalVisible, setShowModal] = useState(false)

    const showModal = () => { setShowModal(true) }
    const hideModal = () => { setShowModal(false) }

    const handleNewCardSubmit = (newCardData) => {
        return post('/api/card', newCardData)
            .then(data => {
                setCards(data)
                hideModal()
            })
            .catch(e => {
                throw e.message
            })
    }

    useEffect(() => {
        get('/api/cards')
            .then((cards) => setCards(cards))
    }, [setCards])

    return (
        <div>
            <h1>Cards</h1>
            <input type="button" onClick={showModal} value="New Card" style={{ marginBottom: 25 + 'px' }} />
            <div>
                <Table rowKey="uuid" columns={{
                    'last_four': 'Last Four',
                    'description': 'Description'
                }} rows={cards} />
            </div>
            <Modal headerText="New Card" visible={modalVisible} handleClose={hideModal}>
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
                                if (!/.....*/.test(fieldValue)) {
                                    return `${fieldLabel} must be at least 4 characters long`
                                }
                                return null
                            }
                        }
                    }}
                />
            </Modal>
        </div>
    )
}

export default CardsPage
