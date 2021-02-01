import React, { useState } from 'react'

import LoaderComponent from '../common/LoaderComponent'
import SplitwiseLoginCheck from '../SplitwiseLoginCheck'
import AccountActivitiesTable from '../tables/AccountActivitiesTable'
import CardActivitiesTable from '../tables/CardActivitiesTable'
import SplitwiseExpenseTable from '../tables/SplitwiseExpenseTable'

const SplitwiseExpensePage = ({ match }) => {
    const splitwiseExpenseUUID = match.params.splitwiseExpenseUUID

    const [splitwiseExpense, setSplitwiseExpense] = useState(null)


    let cardActivitiesTable = null
    if (splitwiseExpense?.card_activities !== null && splitwiseExpense?.card_activities.length > 0) {
        cardActivitiesTable = <div>
            <h3>Card Activities</h3>
            <CardActivitiesTable data={splitwiseExpense?.card_activities} hideFilters={true} />
        </div>
    }


    let accountActivitiesTable = null
    if (splitwiseExpense?.accountivities !== null && splitwiseExpense?.account_activities.length > 0) {
        accountActivitiesTable = <div>
            <h3>Account Activities</h3>
            <AccountActivitiesTable data={splitwiseExpense?.account_activities} hideFilters={true} />
        </div>
    }

    return (
        <div>
            <h1>Splitwise Expense {splitwiseExpenseUUID} ({splitwiseExpense?.description})</h1>
            <LoaderComponent
                path={`/api/splitwise_expenses/${splitwiseExpenseUUID}`}
                parentLoading={false}
                setData={setSplitwiseExpense}
            />
            <SplitwiseLoginCheck>
                <SplitwiseExpenseTable data={splitwiseExpense ? [splitwiseExpense] : []} />
                {cardActivitiesTable}
                {accountActivitiesTable}
            </SplitwiseLoginCheck>
        </div>
    )
}

export default SplitwiseExpensePage
