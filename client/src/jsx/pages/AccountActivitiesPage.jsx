import React, { useState } from 'react'

import { postForm } from '../../utils/api'
import Form from '../common/Form'
import LoaderComponent from '../common/LoaderComponent'
import Modal from '../common/Modal'
import AccountActivitiesTable from '../tables/AccountActivitiesTable'

const AccountActivitiesPage = ({ match }) => {
    const accountUUID = match.params.accountUUID

    const [account, setAccount] = useState(null)
    const [modalVisible, setShowModal] = useState(false)
    const [uploading, setUploading] = useState(false)

    const showModal = () => { setShowModal(true) }
    const hideModal = () => { setShowModal(false) }

    const handleActivityUpload = (activityData) => {
        setUploading(true)
        let formData = new FormData()
        // todo: consider adding support for multiple files
        formData.append('file', activityData['file'])
        return postForm(`/api/accounts/${accountUUID}/activities`, formData)
            .then(() => {
                hideModal()
            })
            .catch(e => { throw e.message })
            .finally(() => setUploading(false))
    }

    return (
        <div>
            <h1>Account Activities for {account ? (account.last_four + " (" + account.description + ")") : null}</h1>
            <input type='button' onClick={showModal} value='Upload Activities' style={{ marginBottom: 25 + 'px' }} />
            <LoaderComponent
                path={`/api/accounts/${accountUUID}`}
                parentLoading={uploading}
                setData={(account) => setAccount(account)}
            />
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
