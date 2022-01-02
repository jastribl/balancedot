import React, { useState } from 'react'

import { postForm, postJSON } from '../../utils/api'
import Form from '../common/Form'
import LoaderComponent from '../common/LoaderComponent'
import Modal from '../common/Modal'
import AccountActivitiesTable from '../tables/AccountActivitiesTable'

const AccountActivitiesPage = ({ match }) => {
    const accountUUID = match.params.accountUUID

    const [account, setAccount] = useState(null)
    const [modalVisible, setShowModal] = useState(false)
    const [uploading, setUploading] = useState(false)
    const [isAutoLinking, setIsAutoLinking] = useState(false)

    const showUploadModal = () => { setShowModal(true) }
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

    const handleAutoLinkCards = () => {
        setIsAutoLinking(true)
        return postJSON(`/api/accounts/${accountUUID}/auto_link_with_card_activities`)
            .catch(e => setErrorMessage(e.message))
            .finally(() => setIsAutoLinking(false))
    }
    return (
        <div>
            <h1>Account Activities for {account ? (`${account.last_four} (${account.description})`) : null}</h1>
            <input type='button' onClick={showUploadModal} value='Upload Activities' style={{ marginBottom: 25 + 'px' }} />
            <br />
            <input type='button' onClick={handleAutoLinkCards} value='Auto Link with Card Activities' style={{ marginBottom: 25 + 'px' }} />
            <LoaderComponent
                path={`/api/accounts/${accountUUID}`}
                parentLoading={uploading || isAutoLinking}
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
