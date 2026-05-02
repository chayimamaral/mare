/* eslint-disable react-hooks/exhaustive-deps */
/* eslint-disable @next/next/no-img-element */

import PrimeReact from 'primereact/api';
import { Button } from 'primereact/button';
import { Sidebar } from 'primereact/sidebar';
import { classNames } from 'primereact/utils';
import React, { useContext, useEffect, useState } from 'react';
import { AppConfigProps, LayoutConfig, LayoutState } from '../types/types';
import { LayoutContext } from './context/layoutcontext';

const AppConfig = (props: AppConfigProps) => {
    /** px na raiz (html); faixa 10–14 para texto mais compacto em telas Windows com DPI alto. */
    const [scales] = useState([10, 11, 12, 13, 14]);
    const { layoutConfig, setLayoutConfig, layoutState, setLayoutState } = useContext(LayoutContext);

    const onConfigButtonClick = () => {
        setLayoutState((prevState: LayoutState) => ({ ...prevState, configSidebarVisible: true }));
    };

    const onConfigSidebarHide = () => {
        setLayoutState((prevState: LayoutState) => ({ ...prevState, configSidebarVisible: false }));
    };

    const changeTheme = (theme: string, colorScheme: string) => {
        PrimeReact.changeTheme?.(layoutConfig.theme, theme, 'theme-css', () => {
            setLayoutConfig((prevState: LayoutConfig) => ({ ...prevState, theme, colorScheme }));
        });
    };

    const decrementScale = () => {
        setLayoutConfig((prevState: LayoutConfig) => ({
            ...prevState,
            scale: Math.max(scales[0], prevState.scale - 1),
        }));
    };

    const incrementScale = () => {
        setLayoutConfig((prevState: LayoutConfig) => ({
            ...prevState,
            scale: Math.min(scales[scales.length - 1], prevState.scale + 1),
        }));
    };

    const applyScale = () => {
        document.documentElement.style.fontSize = layoutConfig.scale + 'px';
    };

    useEffect(() => {
        applyScale();
    }, [layoutConfig.scale]);

    return (
        <>
            <button className="layout-config-button p-link" type="button" onClick={onConfigButtonClick} aria-label="Configurações">
                <i className="pi pi-cog"></i>
            </button>

            <Sidebar visible={layoutState.configSidebarVisible} onHide={onConfigSidebarHide} position="right" className="layout-config-sidebar w-20rem">
                {!props.simple && (
                    <>
                        <h5>Escala</h5>
                        <div className="flex align-items-center flex-nowrap">
                            <Button
                                icon="pi pi-minus"
                                type="button"
                                onClick={decrementScale}
                                rounded
                                text
                                ripple={false}
                                className="w-2rem h-2rem mr-2 flex-shrink-0"
                                disabled={layoutConfig.scale === scales[0]}
                            />
                            <div className="flex gap-2 align-items-center justify-content-center flex-grow-1 min-w-0" aria-hidden>
                                {scales.map((item) => (
                                    <i
                                        key={item}
                                        className={classNames('pi pi-circle-fill text-xs m-0', {
                                            'text-primary-500': item === layoutConfig.scale,
                                            'text-300': item !== layoutConfig.scale,
                                        })}
                                        style={{ lineHeight: 1 }}
                                    />
                                ))}
                            </div>
                            <Button
                                icon="pi pi-plus"
                                type="button"
                                onClick={incrementScale}
                                rounded
                                text
                                ripple={false}
                                className="w-2rem h-2rem ml-2 flex-shrink-0"
                                disabled={layoutConfig.scale === scales[scales.length - 1]}
                            />
                        </div>

                        <h5>Tipo de Menu</h5>
                        <div className="flex">
                            <div className="layout-config-native-radio flex align-items-center gap-2 flex-1">
                                <input
                                    type="radio"
                                    id="mode1"
                                    name="layout-menu-mode"
                                    className="m-0"
                                    checked={layoutConfig.menuMode === 'static'}
                                    onChange={() => setLayoutConfig((p: LayoutConfig) => ({ ...p, menuMode: 'static' }))}
                                />
                                <label htmlFor="mode1" className="m-0 cursor-pointer">
                                    Estático
                                </label>
                            </div>
                            <div className="layout-config-native-radio flex align-items-center gap-2 flex-1">
                                <input
                                    type="radio"
                                    id="mode2"
                                    name="layout-menu-mode"
                                    className="m-0"
                                    checked={layoutConfig.menuMode === 'overlay'}
                                    onChange={() => setLayoutConfig((p: LayoutConfig) => ({ ...p, menuMode: 'overlay' }))}
                                />
                                <label htmlFor="mode2" className="m-0 cursor-pointer">
                                    Escondido
                                </label>
                            </div>
                        </div>

                        <h5>Tipo de Campos</h5>
                        <div className="flex">
                            <div className="layout-config-native-radio flex align-items-center gap-2 flex-1">
                                <input
                                    type="radio"
                                    id="outlined_input"
                                    name="layout-input-style"
                                    className="m-0"
                                    checked={layoutConfig.inputStyle === 'outlined'}
                                    onChange={() => setLayoutConfig((p: LayoutConfig) => ({ ...p, inputStyle: 'outlined' }))}
                                />
                                <label htmlFor="outlined_input" className="m-0 cursor-pointer">
                                    Vazio
                                </label>
                            </div>
                            <div className="layout-config-native-radio flex align-items-center gap-2 flex-1">
                                <input
                                    type="radio"
                                    id="filled_input"
                                    name="layout-input-style"
                                    className="m-0"
                                    checked={layoutConfig.inputStyle === 'filled'}
                                    onChange={() => setLayoutConfig((p: LayoutConfig) => ({ ...p, inputStyle: 'filled' }))}
                                />
                                <label htmlFor="filled_input" className="m-0 cursor-pointer">
                                    Hachurado
                                </label>
                            </div>
                        </div>

                    </>
                )}

                <h5>PrimeOne Design - 2021</h5>
                <div className="grid">
                    <div className="col-3">
                        <button className="p-link w-2rem h-2rem" onClick={() => changeTheme('saga-blue', 'light')}>
                            <img src="/layout/images/themes/saga-blue.png" className="w-2rem h-2rem" alt="Saga Blue" />
                        </button>
                    </div>
                    <div className="col-3">
                        <button className="p-link w-2rem h-2rem" onClick={() => changeTheme('saga-green', 'light')}>
                            <img src="/layout/images/themes/saga-green.png" className="w-2rem h-2rem" alt="Saga Green" />
                        </button>
                    </div>
                    <div className="col-3">
                        <button className="p-link w-2rem h-2rem" onClick={() => changeTheme('saga-orange', 'light')}>
                            <img src="/layout/images/themes/saga-orange.png" className="w-2rem h-2rem" alt="Saga Orange" />
                        </button>
                    </div>
                    <div className="col-3">
                        <button className="p-link w-2rem h-2rem" onClick={() => changeTheme('saga-purple', 'light')}>
                            <img src="/layout/images/themes/saga-purple.png" className="w-2rem h-2rem" alt="Saga Purple" />
                        </button>
                    </div>
                    <div className="col-3">
                        <button className="p-link w-2rem h-2rem" onClick={() => changeTheme('vela-blue', 'dark')}>
                            <img src="/layout/images/themes/vela-blue.png" className="w-2rem h-2rem" alt="Vela Blue" />
                        </button>
                    </div>
                    <div className="col-3">
                        <button className="p-link w-2rem h-2rem" onClick={() => changeTheme('vela-green', 'dark')}>
                            <img src="/layout/images/themes/vela-green.png" className="w-2rem h-2rem" alt="Vela Green" />
                        </button>
                    </div>
                    <div className="col-3">
                        <button className="p-link w-2rem h-2rem" onClick={() => changeTheme('vela-orange', 'dark')}>
                            <img src="/layout/images/themes/vela-orange.png" className="w-2rem h-2rem" alt="Vela Orange" />
                        </button>
                    </div>
                    <div className="col-3">
                        <button className="p-link w-2rem h-2rem" onClick={() => changeTheme('vela-purple', 'dark')}>
                            <img src="/layout/images/themes/vela-purple.png" className="w-2rem h-2rem" alt="Vela Purple" />
                        </button>
                    </div>
                    <div className="col-3">
                        <button className="p-link w-2rem h-2rem" onClick={() => changeTheme('arya-blue', 'dark')}>
                            <img src="/layout/images/themes/arya-blue.png" className="w-2rem h-2rem" alt="Arya Blue" />
                        </button>
                    </div>
                    <div className="col-3">
                        <button className="p-link w-2rem h-2rem" onClick={() => changeTheme('arya-green', 'dark')}>
                            <img src="/layout/images/themes/arya-green.png" className="w-2rem h-2rem" alt="Arya Green" />
                        </button>
                    </div>
                    <div className="col-3">
                        <button className="p-link w-2rem h-2rem" onClick={() => changeTheme('arya-orange', 'dark')}>
                            <img src="/layout/images/themes/arya-orange.png" className="w-2rem h-2rem" alt="Arya Orange" />
                        </button>
                    </div>
                    <div className="col-3">
                        <button className="p-link w-2rem h-2rem" onClick={() => changeTheme('arya-purple', 'dark')}>
                            <img src="/layout/images/themes/arya-purple.png" className="w-2rem h-2rem" alt="Arya Purple" />
                        </button>
                    </div>
                </div>
            </Sidebar>
        </>
    );
};

export default AppConfig;
