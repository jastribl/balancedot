import React, { useEffect, useState } from 'react'

import { get, postForm } from '../../utils/api'
import Form from '../common/Form'
import Modal from '../common/Modal'
import AccountActivitiesTable from '../tables/AccountActivitiesTable'

const AccountActivitiesPage = ({ match }) => {
    const accountUUID = match.params.accountUUID

    const [account, setAccount] = useState(null)
    const [accountActivities, setAccountActivities] = useState([])
    const [modalVisible, setShowModal] = useState(false)

    const showModal = () => { setShowModal(true) }
    const hideModal = () => { setShowModal(false) }

    const refreshAccount = () => {
        get(`/api/accounts/${accountUUID}`)
            .then(accountResponse => setAccount(accountResponse))
    }

    const refreshAccountActivities = () => {
        get(`/api/accounts/${accountUUID}/activities`)
            .then(accountActivitiesResponse => setAccountActivities(accountActivitiesResponse))
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
            <h1>Account Activities for {account ? (account.last_four + " (" + account.description + ")") : null}</h1>
            <input type='button' onClick={showModal} value='Upload Activities' style={{ marginBottom: 25 + 'px' }} />
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
