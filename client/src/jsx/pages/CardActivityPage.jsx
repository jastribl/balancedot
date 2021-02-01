import React, { useEffect, useState } from 'react'

import { get } from '../../utils/api'
import CardActivitiesTable from '../tables/CardActivitiesTable'
import SplitwiseExpenseTable from '../tables/SplitwiseExpenseTable'

const CardActivityPage = ({ match }) => {
    const cardUUID = match.params.cardUUID
    const cardActivityUUID = match.params.cardActivityUUID

    const [card, setCard] = useState(null)
    const [cardActivity, setCardActivity] = useState(null)

    const refreshCard = () => {
        get(`/api/cards/${cardUUID}`)
            .then(cardResponse => setCard(cardResponse))
    }

    const refreshCardActivity = () => {
        get(`/api/card_activities/${cardActivityUUID}`)
            .then(cardActivityResponse => setCardActivity(cardActivityResponse))
    }

    useEffect(() => {
        refreshCard()
        refreshCardActivity()
    }, [setCard, setCardActivity])

    let splitwiseExpenseTable = null
    if (cardActivity?.splitwise_expenses !== null && cardActivity?.splitwise_expenses.length > 0) {
        splitwiseExpenseTable = <div>
            <h3>Splitwise Expenses</h3>
            <SplitwiseExpenseTable data={cardActivity?.splitwise_expenses} hideFilters={true} />
        </div>
    }

    return (
        <div>
            <h1>Card Activity {cardActivityUUID} ({cardActivity?.description}) </h1>
            <h2>For card {card ? (card.last_four + " (" + card.description + ")") : null}</h2>
            <CardActivitiesTable
                data={cardActivity ? [cardActivity] : []}
                hideFilters={true}
            />
            {splitwiseExpenseTable}
        </div>
    )
}

export default CardActivityPage
