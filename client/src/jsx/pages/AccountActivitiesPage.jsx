import React, { useEffect, useState } from 'react'

import { get, postForm } from '../../utils/api'
import { formatAsDate, formatAsMoney } from '../../utils/format'
import { dateComparator } from '../../utils/sorting'
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
            .then(accountActivitiesResponse => setAccountActivities(accountActivitiesResponse))
    }

    const handleActivityUpload = (activityData) => {
        let formData = new FormData()
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
            <h1>Account Activities for {account ? (account.last_four + " (" + account.description + ")") : null}</h1>
            <input type='button' onClick={showModal} value='Upload Activities' style={{ marginBottom: 25 + 'px' }} />
            <div>
                <Table
                    rowKey='uuid'
                    rows={accountActivities}
                    columns={['uuid', 'details', 'posting_date', 'description', 'amount', 'type']}
                    customRenders={{
                        'posting_date': (data) => formatAsDate(data['posting_date']),
                        'amount': (data) => formatAsMoney(data['amount']),
                    }}
                    initialSortColumn='posting_date'
                    customSortComparators={{
                        'posting_date': dateComparator
                    }}
                />
            </div>
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
