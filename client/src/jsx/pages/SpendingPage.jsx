import Moment from 'moment'
import React, { useEffect, useState } from 'react'

import LoaderComponent from '../common/LoaderComponent'

const SpendingPage = ({ match }) => {

    const [accountActivities, setAccountActivities] = useState(null)
    const [cardActivities, setCardActivities] = useState(null)
    const [splitwiseExpenses, setSplitwiseExpenses] = useState(null)

    const spendAccountActivities = accountActivities?.filter(accountActivity => {
        // Filter out payments (money being added to the account)
        // todo: need to consider returns and being paid back for things maybe (although should be on splitwise)
        if (accountActivity.amount >= 0) {
            return false
        }

        // Filter out payments to credit cards (we'll cover those later)
        if (accountActivity.card_activities.length > 0) {
            const amountOnCards = accountActivity.card_activities
                .map(cardActivity => cardActivity.amount)
                .reduce((a, b) => a + b)
            if (-accountActivity.amount === amountOnCards) {
                return false
            }
        }

        const removeDescriptionContains = [
            // Filter out ATM withdrawls for now
            // TODO: Add support for ATM withdrawls
            "WITHDRAWAL",
            "BKOFAMERICA ATM",
            "WITHDRWL",
            "Customer Withdrawal Image",
            // Filter out transfers to Schwab
            "SCHWAB BROKERAGE DES:MONEYLINK",
            "SCHWAB BROKERAGE MONEYLINK",
        ]
        if (removeDescriptionContains.filter(descriptionContains =>
            accountActivity.description.includes(descriptionContains)
        ).length > 0) {
            return false
        }

        if (accountActivity.description.includes("VENMO")) {
            return false
        }

        // Default to keeping it
        return true
    })
    const spendCardActivities = cardActivities?.filter(cardActivity => {
        // Filter out returns
        // TODO: figure out how to deal with returns at some point
        if (cardActivity.type === "Return") {
            return false
        }

        // Filter out adjustments
        // TODO: Figure out how to deal with adjustments at some point
        if (cardActivity.type === "Adjustment") {
            return false
        }

        // Filter out old things from BofA Card since we don't have proper history
        if (cardActivity.amount > 0 &&
            cardActivity.description.includes("Online payment from CHK") &&
            cardActivity.type.includes("Account: 4931") &&
            Moment.utc(cardActivity.post_date) < Moment("2019/05/26")) {
            return false
        }

        // Filter out payments
        // todo: need to consider returns and being paid back for things maybe (although should be on splitwise)
        if (cardActivity.amount > 0 && cardActivity.account_activities.length > 0) {
            const amountOnAccounts = cardActivity.account_activities
                .map(accountActivity => accountActivity.amount)
                .reduce((a, b) => a + b)
            if (-amountOnAccounts === cardActivity.amount) {
                return false
            }
        }

        // Default to keeping it
        return true
    })
    const spendSplitwiseExpenses = splitwiseExpenses?.filter(splitwiseExpense => {
        // Filter out venmo payments
        // TODO: Figure out if this is legit to filter out
        if (splitwiseExpense.description === "Payment") {
            if (splitwiseExpense.creation_method === "venmo") {
                return false
            }
        }

        // Filter out one-to-one debt simplification
        // TODO: Figure out if this is legit to filter out
        if (splitwiseExpense.description === "One-on-one debt simplification") {
            if (splitwiseExpense.creation_method === "debt_consolidation") {
                return false
            }
        }

        // Default to keeping it
        return true
    })

    let dateToExpenses = {}
    spendAccountActivities?.forEach(accountActivity => {
        const date = accountActivity.posting_date.substring(0, 10)
        dateToExpenses[date] ??= {}
        dateToExpenses[date]['account'] ??= []
        dateToExpenses[date]['account'].push(accountActivity)
    });
    spendCardActivities?.forEach(cardActivity => {
        const date = cardActivity.post_date.substring(0, 10)
        dateToExpenses[date] ??= {}
        dateToExpenses[date]['card'] ??= []
        dateToExpenses[date]['card'].push(cardActivity)
    })
    spendSplitwiseExpenses?.forEach(splitwiseExpense => {
        const date = splitwiseExpense.date.substring(0, 10)
        dateToExpenses[date] ??= {}
        dateToExpenses[date]['splitwise'] ??= []
        dateToExpenses[date]['splitwise'].push(splitwiseExpense)
    })
    return (
        <div>
            <h1>Spending Breakdown</h1>
            {/* <pre>{JSON.stringify(spendAccountActivities, null, 4)}</pre>
            <pre>{JSON.stringify(spendCardActivities, null, 4)}</pre>
            <pre>{JSON.stringify(spendSplitwiseExpenses, null, 4)}</pre> */}
            <pre>{JSON.stringify(dateToExpenses, null, 4)}</pre>
            <LoaderComponent
                path={`/api/account_activities/`}
                setData={(accountActivities) => setAccountActivities(accountActivities)}
            />
            <LoaderComponent
                path={`/api/card_activities/`}
                setData={(cardActivities) => setCardActivities(cardActivities)}
            />
            <LoaderComponent
                path={`/api/splitwise_expenses`}
                setData={(splitwiseExpenses) => setSplitwiseExpenses(splitwiseExpenses)}
            />
        </div>
    )
}

export default SpendingPage