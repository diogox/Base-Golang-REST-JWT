import React, { Component } from 'react'
import { withRouter } from 'react-router-dom'
import { RouteComponentProps } from "react-router"
import { Link } from 'react-router-dom'
import { logout } from '../../utils/AuthService'

interface Route {
    title: string,
    route?: string,
    action?: () => void,
}

type NavProps = RouteComponentProps & {
    isLoggedIn: boolean,
    location: Location,
}

class NavBar extends Component<NavProps> {
    constructor(props: NavProps) {
        super(props)
    }

    authedRoutes: Route[] = [
        {
            title: "My Dashboard",
            route: "/dashboard",
        },
        {
            title: "Logout",
            action: () => {
                logout()
                this.props.history.push('/')
            },
        },
    ]
    
    unAuthedRoutes: Route[] = [
        {
            title: "Login",
            route: "/login",
        },
        {
            title: "Signup",
            route: "/signup",
        },
    ]

    
    _renderLink = (title: string, route?: string, action?: () => void) => {
        // Check if it's the currently selected one
        let isActiveStyle = this.props.location.pathname === route? 'text-grey-darker' : 'text-grey'
        let linkStyle = "block no-underline mt-4 lg:inline-block lg:mt-0 hover:text-grey-darker ml-4 cursor-pointer " + isActiveStyle
        
        // If there's an action, it's not linking anywhere
        if (action !== undefined) {
            return (
                <a onClick={action} className={linkStyle}>
                        { title }
                </a>
            )
        }
        
        // Check if route exists
        if (route === undefined) {
            return
        }

        // Link to another route
        return (
            <Link to={route} className={linkStyle}>
                { title }
            </Link>
        )
    }

    render() {
        return (
            <nav className="flex items-center justify-between flex-wrap p-5 mt-2">
                <div className="flex items-center flex-no-shrink text-white mr-6 ml-6">
                    <Link to="/" style={{ textDecoration: 'none' }}>
                        <span className="font-semibold text-xl tracking-tight text-grey-darkest">Logo</span>
                    </Link>
                </div>
                <div className="block lg:hidden">
                    <button className="flex items-center px-3 py-2 border rounded text-teal-lighter border-teal-light hover:text-white hover:border-white">
                        <svg className="fill-current h-3 w-3" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><title>Menu</title><path d="M0 3h20v2H0V3zm0 6h20v2H0V9zm0 6h20v2H0v-2z"/></svg>
                    </button>
                </div>
                <div className="w-full block flex-grow lg:flex lg:items-center lg:w-auto">
                    <div className="text-sm lg:flex-grow text-right mr-6">
                        {
                            this.props.isLoggedIn ?
                                this.authedRoutes.map((obj: Route) => this._renderLink(obj.title, obj.route, obj.action))
                                : this.unAuthedRoutes.map((obj: Route) => this._renderLink(obj.title, obj.route, obj.action))
                        }
                        {/*
                            <Link to="/about">
                                <a href="#responsive-header" className="block mt-4 lg:inline-block lg:mt-0 text-grey hover:text-grey-darker">
                                    About
                                </a>
                            </Link>
                        */}
                    </div>
                    {/*
                        <div>
                            <a href="#" className="inline-block text-sm px-4 py-2 leading-none border rounded text-teal border-teal hover:border-transparent hover:text-white hover:bg-teal mt-4 lg:mt-0">Download</a>
                        </div>
                    */}
                </div>
            </nav>
        )
    }
}

export default withRouter(NavBar)