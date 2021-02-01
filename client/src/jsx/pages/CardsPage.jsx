import React, { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'

import { get, postJSON } from '../../utils/api'
import Form from '../common/Form'
import Modal from '../common/Modal'
import Table from '../common/Table'

const CardsPage = () => {
    const [cards, setCards] = useState(null)
    const [modalVisible, setShowModal] = useState(false)

    const showModal = () => { setShowModal(true) }
    const hideModal = () => { setShowModal(false) }

    const refreshCards = () => {
        get('/api/cards')
            .then(cardsResponse => setCards(cardsResponse))
    }

    const handleNewCardSubmit = (newCardData) => {
        return postJSON('/api/card', newCardData)
            .then(() => {
                hideModal()
                refreshCards()
            })
            .catch(e => {
                throw e.message
            })
    }

    useEffect(() => {
        refreshCards()
    }, [setCards])

    return (
        <div>
            <h1>Cards</h1>
            <input type='button' onClick={showModal} value='New Card' style={{ marginBottom: 25 + 'px' }} />
            <Table
                rowKey='uuid'
                rows={cards}
                columns={['last_four', 'description', 'bank_name']}
                customRenders={{
                    'last_four': (data) =>
                        <Link to={'/cards/' + data['uuid'] + '/activities'}>{data['last_four']}</Link>
                }}
            />
            <Modal headerText='New Card' visible={modalVisible} handleClose={hideModal}>
                <Form
                    onSubmit={handleNewCardSubmit}
                    fieldInfos={{
                        last_four: {
                            inputType: 'text',
                            validate: (fieldLabel, fieldValue) => {
                                if (!/^[0-9][0-9][0-9][0-9]$/.test(fieldValue)) {
                                    return `${fieldLabel} must follow '####' format`
                                }
                                return null
                            }
                        },
                        description: {
                            inputType: 'text',
                            validate: (fieldLabel, fieldValue) => {
                                if (!/.....*/.test(fieldValue)) {
                                    return `${fieldLabel} must be at least 4 characters long`
                                }
                                return null
                            }
                        },
                        bank_name: {
                            inputType: 'select',
                            selectOptions: ['chase', 'bofa'],
                        }
                    }}
                />
            </Modal>
        </div>
    )
}

export default CardsPage
