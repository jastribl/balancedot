import React from 'react'
import { BrowserRouter, NavLink, Route, Switch } from 'react-router-dom'

import AccountActivitiesPage from './pages/AccountActivitiesPage'
import AccountsPage from './pages/AccountsPage'
import CardActivitiesPage from './pages/CardActivitiesPage'
import CardsPage from './pages/CardsPage'
import ErrorPage from './pages/ErrorPage'
import HomePage from './pages/HomePage'
import LinkerFlowPage from './pages/LinkerFlowPage'
import LinkerPage from './pages/LinkerPage'
import OauthCallbackPage from './pages/OauthCallbackPage'
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
                    <NavLink className='nav-item' activeClassName='active-nav-item' to='/linker'>Linker</NavLink>
                </div>
                <Switch>
                    <Route path='/' component={HomePage} exact />
                    <Route path='/accounts' component={AccountsPage} exact />
                    <Route path='/accounts/:accountUUID/activities' component={AccountActivitiesPage} />
                    <Route path='/cards' component={CardsPage} exact />
                    <Route path='/cards/:cardUUID/activities' component={CardActivitiesPage} />
                    <Route path='/splitwise_expenses' component={SplitwiseExpensesPage} />
                    <Route path='/oauth_callback' component={OauthCallbackPage} />
                    <Route path='/linker' component={LinkerPage} exact />
                    <Route path='/linker/:splitwiseExpenseUUID' component={LinkerFlowPage} />
                    <Route component={ErrorPage} />
                </Switch>
            </div>
        </BrowserRouter>
    </div >
)

export default App
