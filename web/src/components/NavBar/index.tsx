import React, { Component } from 'react'
import { Link } from 'react-router-dom';
import { logout } from '../../utils/AuthService'

interface Route {
    title: string,
    route?: string,
    action?: () => void,
}

const authedRoutes: Route[] = [
    {
        title: "My Dashboard",
        route: "/dashboard",
    },
    {
        title: "Logout",
        route: "/",
        action: () => logout(),
    },
]

const unAuthedRoutes: Route[] = [
    {
        title: "Login",
        route: "/login",
    },
    {
        title: "Signup",
        route: "/signup",
    },
]

const LinkMaker = (title: string, route?: string, action?: () => void) => {
    if (route === undefined) {
        return
    }

    return (
        <Link to={route}>
            {
                action !== undefined?
                <a onClick={action}
                   className="block mt-4 lg:inline-block lg:mt-0 text-grey hover:text-grey-darker ml-4 cursor-pointer">
                        { title }
                </a>
                : <a className="block mt-4 lg:inline-block lg:mt-0 text-grey hover:text-grey-darker ml-4 cursor-pointer">
                      { title }
                  </a>
            }
            
        </Link>
    )
}

interface NavProps {
    isLoggedIn: boolean,
    location: Location,
}

export default class NavBar extends Component<NavProps> {
    constructor(props: NavProps) {
        super(props)
    }

    render() {
        return (
            <nav className="flex items-center justify-between flex-wrap p-5 mt-2">
                <div className="flex items-center flex-no-shrink text-white mr-6 ml-6">
                    <Link to="/" style={{ textDecoration: 'none' }}>
                        <span className="font-semibold text-xl tracking-tight text-grey-darkest">Logo  - {console.log(this.props.location)}</span>
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
                                authedRoutes.map((obj: Route) => LinkMaker(obj.title, obj.route, obj.action))
                                : unAuthedRoutes.map((obj: Route) => LinkMaker(obj.title, obj.route, obj.action))
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