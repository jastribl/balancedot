import React, { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'

import { getWithHandling, postJSON } from '../../utils/api'
import ErrorRow from '../common/ErrorRow'
import Form from '../common/Form'
import Modal from '../common/Modal'
import Spinner from '../common/Spinner'
import Table from '../common/Table'

const CardsPage = () => {
    const [cards, setCards] = useState(null)
    const [cardsLoading, setCardsLoading] = useState(false)
    const [modalVisible, setShowModal] = useState(false)
    const [errorMessage, setErrorMessage] = useState(null)

    const showModal = () => { setShowModal(true) }
    const hideModal = () => { setShowModal(false) }

    const refreshCards = () => getWithHandling(
        '/api/cards',
        setCards,
        setErrorMessage,
        setCardsLoading
    )

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
    }, [
        setCards,
        setErrorMessage,
        setCardsLoading,
    ])

    return (
        <div>
            <Spinner visible={cardsLoading} />
            <h1>Cards</h1>
            <input type='button' onClick={showModal} value='New Card' style={{ marginBottom: 25 + 'px' }} />
            <ErrorRow message={errorMessage} />
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
