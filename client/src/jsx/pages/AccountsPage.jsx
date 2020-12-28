import React, { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'

import { get, postJSON } from '../../utils/api'
import Form from '../common/Form'
import Modal from '../common/Modal'
import Table from '../common/Table'

const AccountsPage = () => {
    const [accounts, setAccounts] = useState(null)
    const [modalVisible, setShowModal] = useState(false)

    const showModal = () => { setShowModal(true) }
    const hideModal = () => { setShowModal(false) }

    const refreshAccounts = () => {
        get('/api/accounts')
            .then(accountsResponse => setAccounts(accountsResponse))
    }

    const handleNewAccountSubmit = (newAccountData) => {
        return postJSON('/api/account', newAccountData)
            .then(() => {
                hideModal()
                refreshAccounts()
            })
            .catch(e => {
                throw e.message
            })
    }

    useEffect(() => {
        refreshAccounts()
    }, [setAccounts])

    return (
        <div>
            <h1>Accounts</h1>
            <input type='button' onClick={showModal} value='New Account' style={{ marginBottom: 25 + 'px' }} />
            <div>
                <Table rowKey='uuid' columns={{
                    'last_four': 'Last Four',
                    'description': 'Description'
                }} rows={accounts} customRenders={{
                    'last_four': (data) =>
                        <Link to={'/accounts/' + data['uuid'] + '/activities'}>{data['last_four']}</Link>
                }} />
            </div>
            <Modal headerText='New Account' visible={modalVisible} handleClose={hideModal}>
                <Form
                    onSubmit={handleNewAccountSubmit}
                    fieldInfos={{
                        last_four: {
                            fieldLabel: 'Last Four',
                            fieldName: 'last_four',
                            placeholder: 'Last Four...',
                            inputType: 'text',
                            validate: (fieldLabel, fieldValue) => {
                                if (!/^[0-9][0-9][0-9][0-9]$/.test(fieldValue)) {
                                    return `${fieldLabel} must follow '####' format`
                                }
                                return null
                            }
                        },
                        description: {
                            fieldLabel: 'Description',
                            fieldName: 'description',
                            placeholder: 'Description...',
                            inputType: 'text',
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

export default AccountsPage
