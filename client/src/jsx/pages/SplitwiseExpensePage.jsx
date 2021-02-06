import React, { useState } from 'react'
import { Link } from 'react-router-dom'

import { postJSONWithHandling } from '../../utils/api'
import { formatAsDate, formatAsMoney } from '../../utils/format'
import LoaderComponent from '../common/LoaderComponent'
import AccountActivitiesTable from '../tables/AccountActivitiesTable'
import CardActivitiesTable from '../tables/CardActivitiesTable'
import SplitwiseExpenseTable from '../tables/SplitwiseExpenseTable'

const SplitwiseExpensePage = ({ match }) => {
    const editMode = match.path.endsWith('/edit')

    const splitwiseExpenseUUID = match.params.splitwiseExpenseUUID

    const [splitwiseExpense, setSplitwiseExpense] = useState(null)
    const [linking, setLinking] = useState(false)
    const [errorMessage, setErrorMessage] = useState(null)

    const handleLinking = (entity, action, uuid) =>
        postJSONWithHandling(
            `/api/${entity}/${uuid}/${action}/${splitwiseExpenseUUID}`,
            setErrorMessage,
            setLinking,
        )

    let cardActivitiesTable = null
    if (splitwiseExpense?.card_activities !== null && splitwiseExpense?.card_activities.length > 0) {
        cardActivitiesTable = <div>
            <h3>Card Activities</h3>
            <CardActivitiesTable
                data={splitwiseExpense?.card_activities}
                hideFilters={true}
                extraColumns={['unlink']}
                extraCustomRenders={{
                    'unlink': (data) => <input
                        type='button'
                        onClick={() => handleLinking('card_activities', 'unlink', data['uuid'])}
                        value='Unlink'
                        disabled={!editMode}
                        style={!editMode ? {
                            backgroundColor: 'grey',
                            cursor: 'not-allowed'
                        } : {}}
                    />,
                }}
            />
        </div>
    }

    let accountActivitiesTable = null
    if (splitwiseExpense?.accountivities !== null && splitwiseExpense?.account_activities.length > 0) {
        accountActivitiesTable = <div>
            <h3>Account Activities</h3>
            <AccountActivitiesTable
                data={splitwiseExpense?.account_activities}
                hideFilters={true}
                extraColumns={['unlink']}
                extraCustomRenders={{
                    'unlink': (data) => <input
                        type='button'
                        onClick={() => handleLinking('account_activities', 'unlink', data['uuid'])}
                        value='Unlink'
                        disabled={!editMode}
                        style={!editMode ? {
                            backgroundColor: 'grey',
                            cursor: 'not-allowed'
                        } : {}}
                    />,
                }}
            />
        </div>
    }

    let linksDiv = null
    if (splitwiseExpense !== null && (
        splitwiseExpense.card_activity_links || splitwiseExpense.account_activity_links
    )) {
        let cardLinksDiv = null
        if (splitwiseExpense.card_activity_links?.length > 0) {
            cardLinksDiv = <div>
                <h3>Card Activity Links</h3>
                <CardActivitiesTable
                    data={splitwiseExpense.card_activity_links}
                    hideFilters={true}
                    extraColumns={['link']}
                    extraCustomRenders={{
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
                            onClick={() => handleLinking('card_activities', 'link', data['uuid'])}
                            value='Link'
                        />,
                    }}
                />
            </div>
        }
        let accountLinksDiv = null
        if (splitwiseExpense.account_activity_links?.length > 0) {
            accountLinksDiv = <div>
                <h3>Account Activity Links</h3>
                <AccountActivitiesTable
                    data={splitwiseExpense.account_activity_links}
                    hideFilters={true}
                    extraColumns={['link']}
                    extraCustomRenders={{
                        'posting_date': (data) => <div style={{
                            color: formatAsDate(data['posting_date']) === formatAsDate(splitwiseExpense['date']) ? 'green' : null
                        }}>{formatAsDate(data['posting_date'])}</div>,
                        'amount': (data) => <div style={{
                            color: Math.abs(data['amount']) === Math.abs(splitwiseExpense['amount_paid']) ? 'green' : null
                        }}>{formatAsMoney(data['amount'])}</div>,
                        'link': (data) => <input
                            type='button'
                            onClick={() => handleLinking('account_activities', 'link', data['uuid'])}
                            value='Link'
                        />
                    }}
                />
            </div>
        }
        if (cardLinksDiv !== null || accountLinksDiv !== null) {
            linksDiv = <div>
                <h2>Possible Links</h2>
                {cardLinksDiv}
                {accountLinksDiv}
            </div>
        } else {
            linksDiv = <h2>No Links Found ;(</h2>
        }
    }

    return (
        <div>
            <h1>Splitwise Expense {splitwiseExpenseUUID}</h1>
            <LoaderComponent
                path={editMode ?
                    `/api/splitwise_expenses/${splitwiseExpenseUUID}/for_linking` :
                    `/api/splitwise_expenses/${splitwiseExpenseUUID}`}
                parentLoading={linking}
                parentErrorMessage={errorMessage}
                setData={setSplitwiseExpense}
            />
            <SplitwiseExpenseTable
                data={splitwiseExpense ? [splitwiseExpense] : []}
                hideFilters={true}
            />
            {cardActivitiesTable}
            {accountActivitiesTable}

            <Link to={`/splitwise_expenses/${splitwiseExpenseUUID}` + (editMode ? '/' : '/edit')}>
                <input
                    type='button'
                    value={editMode ? 'View' : 'Edit'}
                    style={{ marginTop: 25 + 'px' }}
                />
            </Link>
            {linksDiv}
        </div>
    )
}

export default SplitwiseExpensePage
