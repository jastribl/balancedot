import React from 'react'
import { BrowserRouter, Route, Switch, NavLink } from 'react-router-dom'

import HomePage from './pages/HomePage'
import ErrorPage from './pages/ErrorPage'
import CardsPage from './pages/CardsPage'
import CardActivitiesPage from './pages/CardActivitiesPage'

const App = () => (
    <div id="app">
        <BrowserRouter>
            <div>
                <div className="nav-hold">
                    <NavLink className="nav-item" activeClassName="active-nav-item" to="/" exact>Home</NavLink>
                    <NavLink className="nav-item" activeClassName="active-nav-item" to="/cards">Cards</NavLink>
                </div>
                <Switch>
                    <Route path="/" component={HomePage} exact />
                    <Route path="/cards" component={CardsPage} exact />
                    <Route path="/cards/:cardUUID/activities" component={CardActivitiesPage} />
                    <Route component={ErrorPage} />
                </Switch>
            </div>
        </BrowserRouter>
    </div >
)
export default App
