import React, { useEffect, useState } from 'react'

import Table from "../common/Table"
import Modal from "../common/Modal"
import Form from "../common/Form"

const CardsPage = () => {
    const [cards, setCards] = useState(null)
    const [modalVisible, setShowModal] = useState(false)
    const [modalSubmitting, setModalSubmitting] = useState(false)
    // const [formState, setFormState] = useState({
    //     last_four: null,
    //     description: null,
    // })

    const showModal = () => { setShowModal(true) }
    const hideModal = () => { setShowModal(false) }

    const handleNewCardSubmit = (newCardData) => {
        setModalSubmitting(true)
        return fetch('/api/card', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(newCardData)
        })
            .then(response => response.json())
            .then(data => {
                setCards(data)
                setModalSubmitting(false)
                hideModal()
            });
    }

    useEffect(() => {
        fetch(`/api/cards`)
            .then((res) => res.json())
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
                    disableForm={modalSubmitting}
                    fieldInfos={{
                        last_four: {
                            fieldLabel: "Last Four",
                            fieldName: "last_four",
                            placeholder: "Last Four...",
                            inputType: "text",
                            validate: (fieldLabel, fieldValue) => {
                                if (!/^[0-9][0-9][0-9][0-9]$/.test(fieldValue)) {
                                    return `${fieldLabel} invalid. Must follow '####' format`
                                }
                                return null;
                            }
                        },
                        description: {
                            fieldLabel: "Description",
                            fieldName: "description",
                            placeholder: "Description...",
                            inputType: "text",
                        }
                    }}
                />
            </Modal>
        </div>
    )
}

export default CardsPage
