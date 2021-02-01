import React, { useEffect, useState } from 'react'

import { get, postForm } from '../../utils/api'
import ErrorRow from '../common/ErrorRow'
import Form from '../common/Form'
import Modal from '../common/Modal'
import Spinner from '../common/Spinner'
import AccountActivitiesTable from '../tables/AccountActivitiesTable'

const AccountActivitiesPage = ({ match }) => {
    const accountUUID = match.params.accountUUID

    const [account, setAccount] = useState(null)
    const [accountActivities, setAccountActivities] = useState([])
    const [accountsLoading, setAccountsLoading] = useState(false)
    const [accountActivitiesLoading, setAccountActivitiesLoading] = useState(false)
    const [modalVisible, setShowModal] = useState(false)
    const [errorMessage, setErrorMessage] = useState(null)

    const showModal = () => { setShowModal(true) }
    const hideModal = () => { setShowModal(false) }

    const refreshAccount = () => {
        setAccountsLoading(true)
        get(`/api/accounts/${accountUUID}`)
            .then(accountResponse => setAccount(accountResponse))
            .catch(e => setErrorMessage(e.message))
            .finally(() => setAccountsLoading(false))
    }

    const refreshAccountActivities = () => {
        setAccountActivitiesLoading(true)
        get(`/api/accounts/${accountUUID}/activities`)
            .then(accountActivitiesResponse => setAccountActivities(accountActivitiesResponse))
            .catch(e => setErrorMessage(e.message))
            .finally(() => setAccountActivitiesLoading(false))
    }

    const handleActivityUpload = (activityData) => {
        let formData = new FormData()
        // todo: consider adding support for multiple files
        formData.append('file', activityData['file'])
        return postForm(`/api/accounts/${accountUUID}/activities`, formData)
            .then(() => {
                hideModal()
                refreshAccountActivities()
            })
            .catch(e => {
                throw e.message
            })
    }

    useEffect(() => {
        refreshAccount()
        refreshAccountActivities()
    }, [setAccount, setAccountActivities])

    return (
        <div>
            <Spinner visible={accountsLoading || accountActivitiesLoading} />
            <h1>Account Activities for {account ? (account.last_four + " (" + account.description + ")") : null}</h1>
            <input type='button' onClick={showModal} value='Upload Activities' style={{ marginBottom: 25 + 'px' }} />
            <ErrorRow message={errorMessage} />
            <AccountActivitiesTable data={accountActivities} />
            <Modal headerText='Activity Upload' visible={modalVisible} handleClose={hideModal}>
                <Form
                    onSubmit={handleActivityUpload}
                    fieldInfos={{
                        file: {
                            inputType: 'file',
                        },
                    }}
                />
            </Modal>
        </div>
    )
}

export default AccountActivitiesPage
