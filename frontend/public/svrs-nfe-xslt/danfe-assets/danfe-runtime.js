/* Vecontab: substitui jQuery 1.2.6 + nfe-vis.js remotos da SEFAZ (timeout em muitas redes). */
(function () {
    function jQuery(sel) {
        if (typeof sel === 'function') {
            if (document.readyState !== 'loading') {
                sel();
            } else {
                document.addEventListener('DOMContentLoaded', sel);
            }
            return;
        }
        if (sel === document) {
            return {
                ready: function (fn) {
                    if (document.readyState !== 'loading') {
                        fn();
                    } else {
                        document.addEventListener('DOMContentLoaded', fn);
                    }
                },
            };
        }
        var nodes = [];
        if (typeof sel === 'string') {
            try {
                nodes = Array.prototype.slice.call(document.querySelectorAll(sel));
            } catch (e) {
                /* ignore */
            }
        }
        return {
            click: function (handler) {
                for (var i = 0; i < nodes.length; i++) {
                    (function (el) {
                        el.addEventListener('click', function (ev) {
                            handler.call(el, ev);
                        });
                    })(nodes[i]);
                }
            },
        };
    }
    window.jQuery = window.$ = jQuery;

    window.mostraAba = function (n) {
        var nodes = document.querySelectorAll('[id^="aba_nft_"]');
        for (var i = 0; i < nodes.length; i++) {
            nodes[i].style.display = 'none';
        }
        var pane = document.getElementById('aba_nft_' + n);
        if (pane) {
            pane.style.display = 'block';
        }
        var tabs = document.querySelectorAll('#botoes_nft li[id^="tab_"]');
        for (var j = 0; j < tabs.length; j++) {
            tabs[j].className = tabs[j].className.replace(/\bnftselected\b/g, '').replace(/\s+/g, ' ').trim();
        }
        var tab = document.getElementById('tab_' + n);
        if (tab) {
            tab.className = (tab.className + ' nftselected').trim();
        }
    };

    window.EventosEnum = {
        CCE: 1,
        CANC: 2,
        CONF_DEST: 3,
        CIENCIA_OP_DEST: 4,
        DESC_OP_DEST: 5,
        OP_NREALIZADA: 6,
        CTE_AUT: 7,
        CANC_CTE_AUT: 8,
        VIST_SUFRAMA: 9,
        INT_SUFRAMA: 10,
        REG_PAS: 11,
        REG_PAS_BRID: 12,
        CANC_REG_PAS: 13,
        MDFE_AUT: 14,
        MDFE_CANC: 15,
        EPEC: 16,
    };

    window.visualizaEvento = function (_id, _tipo) {
        try {
            window.mostraAba(9);
        } catch (e) {
            /* ignore */
        }
    };
})();
