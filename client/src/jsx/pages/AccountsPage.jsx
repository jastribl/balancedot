import React, { useState } from 'react'
import { Link } from 'react-router-dom'

import { postJSON } from '../../utils/api'
import Form from '../common/Form'
import LoaderComponent from '../common/LoaderComponent'
import Modal from '../common/Modal'
import Table from '../common/Table'

const AccountsPage = () => {
    const [accounts, setAccounts] = useState(null)
    const [modalVisible, setShowModal] = useState(false)
    const [addingNewAccount, setAddingNewAccount] = useState(false)

    const showModal = () => { setShowModal(true) }
    const hideModal = () => { setShowModal(false) }

    const handleNewAccountSubmit = (newAccountData) => {
        setAddingNewAccount(true)
        return postJSON('/api/account', newAccountData)
            .then(() => {
                hideModal()
                refreshAccounts()
            })
            .catch(e => { throw e.message })
            .finally(() => setAddingNewAccount(false))
    }

    return (
        <div>
            <h1>Accounts</h1>
            <input type='button' onClick={showModal} value='New Account' style={{ marginBottom: 25 + 'px' }} />
            <LoaderComponent
                path={'/api/accounts'}
                parentLoading={addingNewAccount}
                setData={setAccounts}
            />
            <Table
                rowKey='uuid'
                rows={accounts}
                columns={['last_four', 'description', 'bank_name']}
                customRenders={{
                    'last_four': (data) =>
                        <Link to={'/accounts/' + data['uuid'] + '/activities'}>{data['last_four']}</Link>
                }}
            />
            <Modal headerText='New Account' visible={modalVisible} handleClose={hideModal}>
                <Form
                    onSubmit={handleNewAccountSubmit}
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

export default AccountsPage
