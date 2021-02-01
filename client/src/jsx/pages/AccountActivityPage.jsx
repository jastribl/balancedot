import React, { useState } from 'react'

import LoaderComponent from '../common/LoaderComponent'
import AccountActivitiesTable from '../tables/AccountActivitiesTable'
import SplitwiseExpenseTable from '../tables/SplitwiseExpenseTable'

const AccountActivityPage = ({ match }) => {
    const accountActivityUUID = match.params.accountActivityUUID

    const [accountActivity, setAccountActivity] = useState(null)

    let splitwiseExpenseTable = null
    if (accountActivity?.splitwise_expenses !== null && accountActivity?.splitwise_expenses.length > 0) {
        splitwiseExpenseTable = <div>
            <h3>Splitwise Expenses</h3>
            <SplitwiseExpenseTable data={accountActivity?.splitwise_expenses} hideFilters={true} />
        </div>
    }

    const account = accountActivity?.account

    return (
        <div>
            <h1>Account Activity {accountActivityUUID} ({accountActivity?.description}) </h1>
            <h2>For account {account ? (account.last_four + " (" + account.description + ")") : null}</h2>
            <LoaderComponent
                path={`/api/account_activities/${accountActivityUUID}`}
                parentLoading={false}
                setData={setAccountActivity}
            />
            <AccountActivitiesTable
                data={accountActivity ? [accountActivity] : []}
                hideFilters={true}
            />
            {splitwiseExpenseTable}
        </div>
    )
}

export default AccountActivityPage
