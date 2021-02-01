import React, { useEffect, useState } from 'react'

import { getWithHandling, postForm } from '../../utils/api'
import ErrorRow from '../common/ErrorRow'
import Form from '../common/Form'
import Modal from '../common/Modal'
import Spinner from '../common/Spinner'
import AccountActivitiesTable from '../tables/AccountActivitiesTable'

const AccountActivitiesPage = ({ match }) => {
    const accountUUID = match.params.accountUUID

    const [account, setAccount] = useState(null)
    const [accountLoading, setAccountLoading] = useState(false)
    const [modalVisible, setShowModal] = useState(false)
    const [errorMessage, setErrorMessage] = useState(null)

    const showModal = () => { setShowModal(true) }
    const hideModal = () => { setShowModal(false) }

    const refreshAccount = () => getWithHandling(
        `/api/accounts/${accountUUID}`,
        setAccount,
        setErrorMessage,
        setAccountLoading
    )

    const handleActivityUpload = (activityData) => {
        let formData = new FormData()
        // todo: consider adding support for multiple files
        formData.append('file', activityData['file'])
        return postForm(`/api/accounts/${accountUUID}/activities`, formData)
            .then(() => {
                hideModal()
                refreshAccount()
            })
            .catch(e => {
                throw e.message
            })
    }

    useEffect(() => {
        refreshAccount()
    }, [
        setAccount,
        setErrorMessage,
        setAccountLoading,
    ])

    return (
        <div>
            <Spinner visible={accountLoading} />
            <h1>Account Activities for {account ? (account.last_four + " (" + account.description + ")") : null}</h1>
            <input type='button' onClick={showModal} value='Upload Activities' style={{ marginBottom: 25 + 'px' }} />
            <ErrorRow message={errorMessage} />
            <AccountActivitiesTable data={account?.activities ?? []} />
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
