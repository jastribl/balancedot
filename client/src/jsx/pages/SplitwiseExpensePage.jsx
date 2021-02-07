import React, { useState } from 'react'
import { Link } from 'react-router-dom'

import { postJSONWithHandling } from '../../utils/api'
import LoaderComponent from '../common/LoaderComponent'
import AccountActivitiesTable from '../tables/AccountActivitiesTable'
import CardActivitiesTable from '../tables/CardActivitiesTable'
import SplitwiseExpenseTable from '../tables/SplitwiseExpenseTable'
import SplitwiseLinkingSection from './sections/SplitwiseLinkingSection'

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
            <SplitwiseLinkingSection
                splitwiseExpense={splitwiseExpense}
                handleLinking={handleLinking}
            />
        </div>
    )
}

export default SplitwiseExpensePage
