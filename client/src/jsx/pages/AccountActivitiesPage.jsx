import React, { useEffect, useState } from 'react'

import { get, postForm } from '../../utils/api'
import { formatAsDate, formatAsMoney } from '../../utils/format'
import Form from '../common/Form'
import Modal from '../common/Modal'
import Table from '../common/Table'

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
            .then(accountActivities => setAccountActivities(accountActivities))
    }

    const handleActivityUpload = (activityData) => {
        let formData = new FormData();
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
    }, [setAccountActivities])

    return (
        <div>
            <h1>Account Activities for {account?.last_four}</h1>
            <input type='button' onClick={showModal} value='Upload Activities' style={{ marginBottom: 25 + 'px' }} />
            <div>
                <Table rowKey='uuid' columns={{
                    'uuid': 'Activity UUID',
                    'details': 'Details',
                    'posting_date': 'Post Date',
                    'description': 'Description',
                    'amount': 'Amount',
                    'type': 'Type',
                }} rows={accountActivities} customRenders={{
                    'posting_date': (data) => formatAsDate(data['posting_date']),
                    'amount': (data) => formatAsMoney(data['amount']),
                }} />
            </div>
            <Modal headerText='Activity Upload' visible={modalVisible} handleClose={hideModal}>
                <Form
                    onSubmit={handleActivityUpload}
                    fieldInfos={{
                        file: {
                            fieldLabel: 'File',
                            fieldName: 'file',
                            placeholder: 'File...',
                            inputType: 'file',
                        },
                    }}
                />
            </Modal>
        </div>
    )
}

export default AccountActivitiesPage
