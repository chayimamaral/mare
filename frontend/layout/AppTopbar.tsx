/* eslint-disable @next/next/no-img-element */

import Link from 'next/link';
import { classNames } from 'primereact/utils';
import React, { forwardRef, useContext, useImperativeHandle, useRef } from 'react';
import { AppTopbarRef } from '../types/types';
import { LayoutContext } from './context/layoutcontext';
import AuthContext from '../components/context/AuthContext';
import { Tooltip } from 'primereact/tooltip';
import { useCaixaPostal } from '../components/context/CaixaPostalContext';

const AppTopbar = forwardRef<AppTopbarRef>((props, ref) => {
    const { layoutState, onMenuToggle, showProfileSidebar } = useContext(LayoutContext);
    const menubuttonRef = useRef(null);
    const topbarmenuRef = useRef(null);
    const topbarmenubuttonRef = useRef(null);
    const { logoutUser, user } = useContext(AuthContext);
    const { naoLidas } = useCaixaPostal();
    const loggedUserName = user?.nome?.trim() || 'Profile';
    const empresaLabel =
        user?.tenant?.nome?.trim() ||
        user?.tenant?.schema_name?.trim() ||
        user?.tenant?.schemaName?.trim() ||
        user?.tenant?.id?.trim() ||
        '';

    useImperativeHandle(ref, () => ({
        menubutton: menubuttonRef.current,
        topbarmenu: topbarmenuRef.current,
        topbarmenubutton: topbarmenubuttonRef.current
    }));

    function handleProfile(): void {
        logoutUser();
    }

    return (
        <div className="layout-topbar">
            <Link href="/" className="layout-topbar-logo">
                {/*}
             <img src={`/layout/images/logo-${layoutConfig.colorScheme !== 'light' ? 'white' : 'dark'}.svg`} width="47.22px" height={'35px'} alt="logo" />
    */}
                <img src="/vecontab.svg" width="47.22px" height={'35px'} alt="Vecontab logo" />

                <span>Vecontab</span>
            </Link>

            <button ref={menubuttonRef} type="button" className="p-link layout-menu-button layout-topbar-button" onClick={onMenuToggle} aria-label="Menu">
                <i className="pi pi-bars" />
            </button>

            <button ref={topbarmenubuttonRef} type="button" className="p-link layout-topbar-menu-button layout-topbar-button" onClick={showProfileSidebar} aria-label="Opções">
                <i className="pi pi-ellipsis-v" />
            </button>

            <div ref={topbarmenuRef} className={classNames('layout-topbar-menu', { 'layout-topbar-menu-mobile-active': layoutState.profileSidebarVisible })}>
                <Link href="/">
                    <Tooltip target=".btn-dashboard" position="bottom" />
                    <button type="button" className="btn-dashboard p-link layout-topbar-button" data-pr-tooltip='Dashboard'>
                        <i className="pi pi-home"></i>
                        <span>Dashboard</span>
                    </button>
                </Link>
                <Link href="/caixa-postal">
                    <Tooltip target=".btn-caixa-postal" position="bottom" />
                    <button type="button" className="btn-caixa-postal p-link layout-topbar-button" data-pr-tooltip="Caixa Postal">
                        <i className={classNames('pi pi-envelope', { 'text-red-500': naoLidas > 0, 'fadein animation-iteration-infinite animation-duration-1000': naoLidas > 0 })}></i>
                        {naoLidas > 0 && (
                            <span className="caixa-postal-badge">{naoLidas > 99 ? '99+' : naoLidas}</span>
                        )}
                        <span>Mensagens</span>
                    </button>
                </Link>
                {empresaLabel && (
                    <>
                        <Tooltip target=".btn-empresa-topbar" position="bottom" />
                        <Link href="/registro">
                            <button type="button" className="btn-empresa-topbar p-link layout-topbar-button layout-topbar-tenant-button" data-pr-tooltip='Dados da Contabilidade'>
                                <i className="pi pi-building"></i>
                                <span>{empresaLabel}</span>
                            </button>
                        </Link>
                    </>
                )}
                <Tooltip target=".btn-login" position="bottom" />
                <button type="button" className="btn-login p-link layout-topbar-button layout-topbar-user-button" data-pr-tooltip='Trocar Usuário' onClick={handleProfile}>
                    <i className="pi pi-user"></i>
                    <span>{loggedUserName}</span>
                </button>
            </div>
        </div>
    );
});

AppTopbar.displayName = 'AppTopbar';

export default AppTopbar;
