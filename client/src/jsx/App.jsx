import React from 'react'
import { BrowserRouter, NavLink, Route, Switch } from 'react-router-dom'

import AccountActivitiesPage from './pages/AccountActivitiesPage'
import AccountActivityPage from './pages/AccountActivityPage'
import AccountsPage from './pages/AccountsPage'
import CardActivitiesPage from './pages/CardActivitiesPage'
import CardActivityPage from './pages/CardActivityPage'
import CardsPage from './pages/CardsPage'
import ErrorPage from './pages/ErrorPage'
import HomePage from './pages/HomePage'
import OauthCallbackPage from './pages/OauthCallbackPage'
import SplitwiseExpensePage from './pages/SplitwiseExpensePage'
import SplitwiseExpensesPage from './pages/SplitwiseExpensesPage'

const App = () => (
    <div id='app'>
        <BrowserRouter>
            <div>
                <div className='nav-hold'>
                    <NavLink className='nav-item' activeClassName='active-nav-item' to='/' exact>Home</NavLink>
                    <NavLink className='nav-item' activeClassName='active-nav-item' to='/accounts'>Accounts</NavLink>
                    <NavLink className='nav-item' activeClassName='active-nav-item' to='/cards'>Cards</NavLink>
                    <NavLink className='nav-item' activeClassName='active-nav-item' to='/splitwise_expenses'>Splitwise Expenses</NavLink>
                </div>
                <Switch>
                    <Route path='/' component={HomePage} exact />
                    <Route path='/accounts' component={AccountsPage} exact />
                    <Route path='/accounts/:accountUUID/activities' component={AccountActivitiesPage} exact />
                    <Route path='/accounts/:accountUUID/activities/:accountActivityUUID' component={AccountActivityPage} />
                    <Route path='/cards' component={CardsPage} exact />
                    <Route path='/cards/:cardUUID/activities' component={CardActivitiesPage} exact />
                    <Route path='/cards/:cardUUID/activities/:cardActivityUUID' component={CardActivityPage} />
                    <Route path='/splitwise_expenses' component={SplitwiseExpensesPage} exact />
                    <Route path='/splitwise_expenses/:splitwiseExpenseUUID' component={SplitwiseExpensePage} />
                    <Route path='/oauth_callback' component={OauthCallbackPage} />
                    <Route component={ErrorPage} />
                </Switch>
            </div>
        </BrowserRouter>
    </div >
)

export default App
