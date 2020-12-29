import Moment from 'moment'
import React, { useEffect, useState } from 'react'

import { get, postJSON } from '../../utils/api'
import { formatAsDate, formatAsMoney } from '../../utils/format'
import ErrorRow from '../common/ErrorRow'
import Spinner from '../common/Spinner'
import Table from '../common/Table'

const LinkerFlowPage = ({ match }) => {
    const splitwiseExpenseUUID = match.params.splitwiseExpenseUUID;

    const [splitwiseExpense, setSplitwiseExpense] = useState(null)
    const [cardsLinks, setCardLinks] = useState(null)
    const [expenseLoading, setExpenseLoading] = useState(false)
    const [cardLinksLoading, setCardLinksLoading] = useState(false)
    const [linkLoading, setLinkLoading] = useState(false)
    const [errorMessage, setErrorMessage] = useState(null)

    useEffect(() => {
        setExpenseLoading(true)
        get(`/api/splitwise_expenses/${splitwiseExpenseUUID}`)
            .then(splitwiseExpenseResponse => {
                setSplitwiseExpense(splitwiseExpenseResponse)
            })
            .catch(e => {
                setErrorMessage(e.message)
            })
            .finally(() => {
                setExpenseLoading(false)
            })

        setCardLinksLoading(true)
        get(`/api/card_activities/for_link/${splitwiseExpenseUUID}`)
            .then(cardLinksResponse => {
                setCardLinks(cardLinksResponse)
            })
            .catch(e => {
                setErrorMessage(e.message)
            })
            .finally(() => {
                setCardLinksLoading(false)
            })
    }, [setExpenseLoading, setSplitwiseExpense, setCardLinks, setErrorMessage])

    const linkCardActivityToExpense = (cardActivityUUID) => {
        setLinkLoading(true)
        postJSON(`/api/card_activities/${cardActivityUUID}/link/${splitwiseExpenseUUID}`)
            .then(data => {
                if (data.message === 'success') {
                    window.history.back()
                }
            })
            .catch(e => {
                setErrorMessage(e.message)
            })
            .finally(() => {
                setLinkLoading(false)
            })
    }

    return (
        <div>
            <Spinner visible={expenseLoading || cardLinksLoading || linkLoading} />
            <h1>Splitwise Expense Linking Flow</h1>
            <ErrorRow message={errorMessage} />
            <div>
                {/* TODO: don't use table for sinle thing */}
                <Table rowKey='uuid' columns={{
                    'uuid': 'UUID',
                    'splitwise_id': 'Splitwise ID',
                    'description': 'Description',
                    'details': 'Details',
                    'amount': 'Amount',
                    'amount_paid': 'Amount Paid',
                    'date': 'Date',
                    'category': 'Category',
                }} rows={splitwiseExpense ? [splitwiseExpense] : []} customRenders={{
                    'details': (data) => data['details'].trim(),
                    'date': (data) =>
                        Moment(data['date']).format('YYYY-MM-DD'),
                    'amount': (data) => formatAsMoney(data['amount'], data['currency_code']),
                    'amount_paid': (data) => formatAsMoney(data['amount_paid'], data['currency_code']),
                }} />
            </div>
            <h2>Possible Links</h2>
            <h3>Card Links</h3>
            <div>
                <Table rowKey='uuid' columns={{
                    'transaction_date': 'Transaction Date',
                    'post_date': 'Post Date',
                    'description': 'Description',
                    'category': 'Category',
                    'type': 'Type',
                    'amount': 'Amount',
                    'link': 'Link',
                }} rows={cardsLinks} customRenders={{
                    'transaction_date': (data) => <div style={{
                        color: formatAsDate(data['transaction_date']) === formatAsDate(splitwiseExpense['date']) ? 'green' : null
                    }}>{formatAsDate(data['transaction_date'])}</div>,
                    'post_date': (data) => <div style={{
                        color: formatAsDate(data['post_date']) === formatAsDate(splitwiseExpense['date']) ? 'green' : null
                    }}>{formatAsDate(data['post_date'])}</div>,
                    'amount': (data) => <div style={{
                        color: Math.abs(data['amount']) === Math.abs(splitwiseExpense['amount_paid']) ? 'green' : null
                    }}>{formatAsMoney(data['amount'])}</div>,
                    'link': (data) => <input
                        type='button'
                        onClick={() => linkCardActivityToExpense(data['uuid'])}
                        value='Link'
                    />
                }} />
            </div>
        </div>
    )
}

export default LinkerFlowPage
