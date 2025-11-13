#!/bin/bash
# install-client.sh
# Script de instalação do Procspy Client para Linux usando systemd
# Requer privilégios de root
#
# Uso:
#   sudo ./install-client.sh [opções]
#
# Opções:
#   -b, --binary PATH       Caminho para o binário procspy-client
#   -c, --config PATH       Caminho para o arquivo de cig-client.json
#   -h, --help              Exibe esta ajuda

set -e

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

# Configurações padrão
SERVICE_NAME="procspy-client"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/procspy"
LOG_DIR="/var/log/procspy"
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"

# Caminhos padrão dos arquivos fonte
DEFAULT_BINARY_PATH="./bin/procspy-client"
DEFAULT_CONFIG_PATH="./etc/config-client.json"

# Variáves customizminhos custopreencs
BINARY_PATH="_PATH=""
CUSTOM_CATH=""

# Função para exibir ajuda
show_he {) 
    cat << EOF
${GREEN}=== Instalação do Procspy Client ===${NC}

Uso: sudo $0 [opções]

Opções:
  -b, --binary PATH       Caminho para o binário procspy-client
                          Padrão: $DEFAULT_BINARY_PATH
  
  -c, --config PATH       Caminho para o arquivo config-cgurat.json
                          Padrão: $DEFAULT_CONFIG_PATH
  
  -h, --help              Exibe esta ajuda

Exemplos:
  # Instalação padrão (binário e config no diretório do projeto)
  sudo $0

  # Instalação com binário customizado
  sudo $0 --binary /tmp/procspy-client

  # Instalação com binário e config customizados
  sudo $0 -b /tmp/procss/procspy -c /tm -c /opt/configsonfig-client.json

  # Instalação a partir de um ixaetório de distribuição
  sudo $0 -b /hot//user/dcspy/procspy-client -nt/mnt/usb/procs/y/config-client.jsoison

Notas:
    Este script deve seexecutatado como root
  - Se os caminhos rspecific especificados, o s atualocura noório atual
  - O binário será copiTALL_DI $INSTALL_DIR
# Parse dfiguração será gumeada para $CONFIG_D
wh- Um serviço systemd sriado e habil
 OF
EOF
    exi
}

# Par-binary)0umentos
while [[ $#]]; do
    arye $1 in
        INARY_PATH="
           2INARY_PATH="$2"
            shi
            ;;
        -c|--config)ONFIG_PATH="$2PATH="$2"
            CONFIG_"$2"
            ;;ift 2
           --help)
        -h|--help)elift 2
            show_help
      ;;
        *)
        echo -Erro: Opção desa: $1${NC}"
    ecse --helpver as opçõeníveis"
         exit 1
            ;;
    e       ;;
done

# Defin   aminhos padrão se não   ram especificados
if       $BINARY_PATH" echthen
    BIN       H="$DEFAULT_BINAR{RED}Er
fi

if [ -z "$ro: OpçãoTH" ]; desc -c|--config${NC}"
    CONFIG_PA     )  echT_CONFIG_Po "Use --er as opções dis"
  

echo -e "${        xit stalação do Procspy Cl1nt ===${NC}"
ho ""
echo -e "${uração:${NC}"
ec"  Binário: $BIN"
echo "  CoCONFIG_PATH"
      "

# Verifica s   stá rodand   omo root
if ;;"$EUID" -ne 0 ]; t
    echo -e "${RED}Erript deve ser exo root${NC}"
  cute: sudo $0"
    exit 1
fi

ec      "${GREEN}✓ Privilég   de root vos${NC}"

a se binárite
if [ ! -f $BINARY_PAhen
    echo -e "rro: Binário não em: $BINARY_PATH$
    echo ""
  ${YELLOW}SoluçõeC}"
 1. Execute o bueiro: ./bui
    echo "  ue o caminho co $0 --bino/para/procst"
    echo "  -help para ver tões"
    exit 1
done

# D  -e "${GREEN}✓   Cário encontrado: ONF caminhos${NC}"
 finaIG_PATH="$2"
 INARY_P   se binário é ex      USTOM_BINARY_PATH   sEFAULT_BINARY_PAThift 2
if [ ! -x "$B     GPATH" ]; then_PATH="${CUSTOM_CONFI   ATH:-$DEFAULT_CO        -h"
echo -e "${YELLso: Binárioé executávelpermissões...${NC}"
  +x "$BINARY_PATH"
ec

# Verificho -e "$figura{Go existe (opcional|-EN}=== Instalação -help)cspy Client      usa"
echo "gXISTS=false
if [ -CONFIG_PAT then
    CONFS=true
    echo -{GREEN}✓ Arquivuração enconCONFIG_PATH${NC
else
    echo YELLOW}⚠ Are configuraçãorado: $CONFI{NC}"
    eche "${YELLOW}  configo padrão será criada${NC
# Verifica serodando como 
     "$EUID    e 0 ]; the    *)
  Para servi      stente
if sys  -etl is-active --q "et $SERVICE_${RED}Errn
    eo:o -e "${YELLOW}Paran Este seço existente..cho "${
    systemctl stove sERVICE_NAME
 er executa
fi

# Copia bináriodRED}Erroção${NC}"
     des"${CYAN}Copiando becidao para $INSTAcu_DIR...${Nte"
cp "$BINARY_: sudo $INSTALL_DIR/pro0"ient"
c +x "$INSTALL_Dpy-clien
 cho -e "${GREEN}✓  : $1$o insta{NC} em $INS"nte...${procspy-clieNC}"CC}"

# Cria diretóE_NA
echo -e "${CME}Criando diretór.${NC}"
fi"$CONFIG_
 sysr -p temctl st"
chmod 755opNSTA_DIR"
echo -e "LL_DST }✓ Diretórios $ALados${NC}"
IR/procs"
echo -a ou cre "${GREguração
if EN}$CONFIG_EXIpyS" = true ]; then
 -cL_Dio -e "${CYAN}Copiando cinsiguração...${NC}talado em $INnt"_DIR/procspy}"
re  cp "$CONFIG_Ptóriosos..FIG_DIR/config-cli.${N
    echo -e GREEConfiguração co para $CONFIG_Dconfig-client.
else
    echo -e"${CYAN}Criando confidrão...${
    cat > NFIG_DIR/config-client.j
{
    "usewhoami)",
   path": "$LOG_DIR
echo"debug": false,
  -p"inte "$LOG 5,
    "serv_DIR"https://seor.com/p
    "api_ht": calhost",
    rt": 8888
}
EO
ch -p ho -e ""-eCY$LOG Configuração pa_DIRANriada em $CONFIG_DG_/confiDIclient.json${R"}"
    echo -e CriYELLando IMPORTANTE: Edite a cond antes de iniciar oNC}"
    e"${YELLOW} e: sudo micro $C/config-client.json"
fi

# Cria de serviço system
# Cr -e "${CYAN}Cri"${GREEquivo de serviço syN}iatórg-ios}"
cat > criados${_FILE" <<EOF
[clientn${NC}"
elseription=Procspy C - Monitto de processos para role paren
    mentationcattps://g >hub.com/se "$o -eIG_Drocspy
After=netIR/conarget
fig- "${CYAN}on" <<EOCfiguração padrão...
}" para $
Type=simple
CONFIG_ot
WorkingDirectoryhoami)",_DIR
ExecStart=$INDIR/coR/procspy-clien$CONFIG_DIR/ct.json
Rert=always
Restec=5
StandardOt=journal
StaardError=journal

# e recurso
    "lOFILE=65536

# Segurogça
NoNe_pat  leges=false
"us$LOG_mp=true

DIR"er":   echo -GREEN}✓ Configuraçãar o serviço$
WantedBy=multi-us    echo 
EOF

echo -e "${GR"${YELLOW} vo de serv E: sudo micm $SERVICE_FILEro $CONFig-client.j}"
fiemdiando arquivo de rviço systemdcspy
A Recarrega systemdfter=network.get
cho -e "${AN}Recarregando sy.${NC}"
systemctleload
echo EEN}✓ Sysecarregado${NC"

# Haba serviço
[o -e "${CYAo serviçoara iniciarboot...${NC
Usetemctl enable $SERVr=roNAME
echo -eot{GREEN}✓ Serviço${NC}"

# Inic
echoAN}Iniciandoviço...${NC}"
systtart $SERVICE_ME
sleep 2
WorkServirectory=$CONFIG_Dice]
Tart=$Iica statusNSTALL_
Rf systemctl iestarDIR/--quiet $SER=5E_NAME; then
  echo -e REEN}✓ Serviçom sucesso${NC}"
StandardOulient $CONFI
Sta echo -e ndardErrG_⚠ Serviço inDIR/nt.js nãodando${NC}"
536-e "${YELLifique osurnalctl -u $SERVICE_50${NC}"
   o -e "${YELLOW}  síveis causas:${NC
ho -e "${YW}    - iguração inválid}"
  o -e "${YELLOrvidor não ace{NC}"
    e "${YELL   - Porta já em 
fi

echo ""
echo -e "${GREEN}=== Instalação Completa ===${NC}"
echo ""
echo -e "${CYAN}Informações do Serviço:${NC}"
echo "  Nome:         $SERVICE_NAME"
echo "  Binário:      $INSTALL_DIR/procspy-client"
echo "  Configuração: $CONFIG_DIR/config-client.json"
echo "  Logs:         $LOG_DIR"
echo "  Service file: $SERVICE_FILE"
echo ""
echo -e "${CYAN}Comandos Úteis:${NC}"
echo "  ${GREEN}Verificar status:${NC}"
echo "    systemctl status $SERVICE_NAME"
echo ""
echo "  ${GREEN}Gerenciar serviço:${NC}"
echo "    sudo systemctl start $SERVICE_NAME       # Iniciar"
echo "    sudo systemctl stop $SERVICE_NAME        # Parar"
echo "    sudo systemctl restart $SERVICE_NAME     # Reiniciar"
echo "    sudo systemctl enable $SERVICE_NAME      # Habilitar auto-start"
echo "    sudo systemctl disable $SERVICE_NAME     # Desabilitar auto-start"
echo ""
echo "  ${GREEN}Ver logs:${NC}"
echo "    journalctl -u $SERVICE_NAME -f           # Tempo real"
echo "    journalctl -u $SERVICE_NAME -n 50        # Últimas 50 linhas"
echo "    journalctl -u $SERVICE_NAME --since today # Logs de hoje"
echo ""
echo "  ${GREEN}Editar configuração:${NC}"
echo "    sudo micro $CONFIG_DIR/config-client.json"
echo "    sudo systemctl restart $SERVICE_NAME     # Após editar"
echo ""
echo -e "${GREEN}✓ Procspy Client instalado com sucesso!${NC}"

# Seguranç
NoNewPrivilege
PrivateTmp=true
[Install]ERVICE_C}"
temdICE_FILE"
echo "
ec{Cçoe "${CYArquivo${NC}"
ech echo "  Configuraço FIG_EXISG_utomaPATH"
fcho "  journalctl -u $ systemctl r -f      # Ver logi # rtempo real"
echo $SERVIrnalctl -u $SERVICECEIME -n 50ni # Ver últimas 50 lcihas de log"
echo ar  udo micro $CONF # Rei/config-client.json   serviçoervifiguração"
ecço$SERVICE_NAME      # Plse ]; then
    earar "
fi
echo -e "$serviçE✓ Procspy Client Lro $CONFLOom sucesso!W}⚠ ATENÇÃO:IG_DIR/o"ite aclient.json conf"
    echigão antes de usaro!${NC}"
    ec -e "${YELLO
ef [ "$CONFIG_EXISTS"lso "  systemctle ""art $SERVIC
echo -e W}Comandos Úteis:ificar sC}"ta
echo "  syste
"  systemctlnstale o Watcher para proteção adicional"
echo ""}✓ Procspy Client instalado com sucesso!${NC}"

echo -e "${GREEN stat_NAME    restart $SERVICE_NAME"
echo "  3. I
  ho "  2. Reinicie o serviço após editar: sudo systemctl  o:ho "  Confi Pasnfiguração se necessário"
ecsosguraçã $B:${NC}"
echo "  1. Edite a coTS" = true INARY_PATH
if [ LOW}Próximos
echo "  Logs: $LOG_DI-Yação"
echo -e "${YELecho ""

echo "  Service file:AN}{CYAN}Habilitando Recvson  # Editar configuriço para iarciar no booren-reload
systemcclient "${GREENNFIG_DIR/config-client.jtl enabiço habilitaleo{NC.json"}"
echo "  sudo micro $CO_DIR/coimas 50 linhas de log"

# "  ConfiguIniração:c$SERVICE_Nando serviço..
firio: $INSTALL_Dpy-client"   # Ver últ
ec
o:${Ne: $SERVICctl -u $SERVICE_NAME -n 50E_C}"NAME"
echo "  journalm tempo real"
echo
echo "E_NAME -f      # Ver logs e
echo "  journalctl -u $SERVIC
echo "e "${GREE# ReiniciaNlaaçõção es do SeComplC}"r serviço
echo -e "${CYAN}o 
systncoremctl start echo -e "${MEreta${NC}
elseE_NAME     
    echo sleep 2i restart $SERVICs-OW}⚠ Sossível causa: confieractive -talado mas não est-quado co${NC}"
    echo -e "${Ym su "${YELLOciet $SERique os logs:VAME; tesso$-u $SERVICE_NA{NChNC}"
  erviçosystemctl"
echo "  
  echo -e ${GREEN}✓ S    # Parar s
#CAM{ifica status
if systeNC}"E_NAME    
ystemctl stop $SERVIC
# Ha
echo -g san${GREEN}✓ Systemd rd${NC}"
syso " temctl o"
ech
# Recarr # Iniciar serviç
Arquivo de serviçcriaE      
echo -e "$stemctl start $SERVICE_NAMWanted# Limiti-user.targetes de reson
Rcho "  syimitNOFILestal status $SERVICE_NAME      # Verificar status"
e
catumentation=htt > "$Sithub.com/seu-usuERVICechoamento de pLE" <<EOFrocontrole pa
[Un -epy ystemctClient - 
De
echo "  sscripti 
}"
# Cria o de servÚteis:${NC
    cpebug": io: $BINARdos Y_PATH"
echo -e "${YELLOW}Coman "  Config: $CONFIG_PATH"
echo ""
echofal "$ilizados:CIR"${NC}"
echo "  Binár
echo "  Service file: $SERVICE_FILE"
echo ""vos Fonte Ut
echo -e "${CYAN}ArquiONF IMPORTANTE: EditeIG_P" "$COração antes_DIR/confi criada g-clienDIR/config-client.jsonNC}"
   o -e "${YELLOW
# Co"server_: "urlol"lhost",
    ": "hport": 8888
}t"in://Logs: $LOG_Dseu-servidor.coteonfiguração padrrocspy",
    echo -e "$    val": 5,
 cho "  pia ou cria chmod +x $CONFIG_DIR/config-client.json"
eIão
    echo -e "$uração: {CYAN}Configuração
echo "  Config
if [ "$CONR."LL_DSTS" = truIR/p
echo "  Binário: $INSTALL_DIR/procspy-client"rocspy
cp "$BINA$rviço: $SERVICE_NAME"{CYAN}Informações da Instalação:${NC}"
echo "  SeRY_PAT
fiho -e "inário
echo
echo -e " -e "${CYA${}"YELLOW}Pinário para $aranço exis
if sctl is-active -ICE_NAME; th
fi""
echo "
_PATH${NC}"nfiguração p criada${a ===${NC}
    ecLLerviço existenOW}  Umomplet
echo  e IN}=== InN"${GRstalação CEELLOW}⚠ Configuração nN}✓ P"
fi"${GREE
e 
echo -echo ""
    rada em:   gis as opções"ARNFIG_PATH${NC}
else journalctl -u $SERVICE_NAME -n 50${NC}
    echo -e os logs:
CONo -eFIG_ "${GEXISTXISTnrifiquefiguração encontradS=trueS=YTH_PAT" ]; then
H}"
if [ cho -e "${YELLOW}  Ve-f "$CONFI${NC}"
    e
# Ve arquivo configuraá rodando
   xit 1talado mas não est
LLOW}⚠ Serviço ins
echo -e     REEN}✓ Binário encecho "  3. Uos d-help pare root vechficados${NC}"
o e --help para ver as s dis"so${NC}"
else"${YE
    echo -e 
    rific""minho/para/procsp
    echo " o buho doild prcomimuild.sh" bri suceso: $0 --binary
  . Especifiqueo iniciado 
    echo a se bGREEN}✓ ServiçYELLOW}Solináes:${NCrio ex   e
if [ ! - -e "${f   itNARY_PATH"  1n${NC}"
   iechove --quiet $SERVICE_NAME; then
    
    echemctl is-act"${RED}Eário não e em: $BINARY_
           cNC}us
if syst"ME
sleep fica stat2

# Veri
viço...${NC}"ERVICE_NA
systemctl start $S
echo -e "${CYAN}Iniciando serviço
# Inicia ser
donene caminhos ão se não forailitado${m ess
if [emctl enable $SERVICE_NhabAME
echo -e "${GREEN}✓ Serviço  -zINARY_PATH" ];BINARY_PULT_BINARY_PATH"
iystf [$CONFIG_PATH" en
s
   ="$DEFAr=jouULrnal
e recursost...${NC}"
LimitNOFILE=65536
...${NC}"regado${NC}"iciar no boo
 "${CYAN}Habilitando serviço para in
echo -e
# Habilita serviço
systemctl${GREEN}✓ Systemd recar daemon-reload
echo -e "
# Segurança
NoNeateTmp=true
o systemd
[Instae "${CYAN}Recarregandll]

echo -
# Recarrega systemdWantedBy=multi-user.target{NC}"

EOREEN}✓ ArquiFvo de serviço criado$
{G
echo -e "$wPrivileges=false
Priv
# Limites dTG_PATH"
finaldErro
Standar
RtandardOutput=jourestartSec=5
Syt $CONFIG_DIR/config-client.json
Restart=always

After=network.targeten

[Service]L_DIR/procspy-cli
Tyer=ngDirectory=$CONFIG_DIR
ExecStart=$INSTALroot
Workipe=simple
Us
GREEN}=== Instalaçãpylio/procsp
Documentation=https://github.com/seu-usuar Client ==="
echo ""ntrole parenta
 para co
# Vefica se estdo como rootssos
if [ "$Ee 0 ]; then
    ech"${REProcspy Client - Monitoramento de procete script deEOF
[Unitiption=]
Descrve so como root$  echo e: sud
    exiecho -e "FILE" <<${GREEN}✓égios de root vecados${NC}"ica se bie
iat > "$SERVICE_f [ ! -f "$TH"{NC}"
c
    echo -e "${rr não en: $..$BINAR}"
    echo"" serviço systemd.
ho -e "${YELLOW}SoluçNC}quivo de"
    stemdCYAN}Criando ar
echo -e "${e"  1. Execute ild primeiro.sh""
o sy
# Cria arquivo de serviç
    ec 2. Eque o inho cor-binary /ca/pent"
  3. Us par$CONFIG_DIR/config-client.jsona ver todas es"
    1
fiecho -e }✓ Binário encot.json${NCntraRY_PAT# Verifica se bin éutável}"
fifiguração
chmod 644 "
da con
# Ajusta permissões 
if [ !NARY_enconfig-clien
    echo -e "}Aviso: B executável, ajustaer{NC}"R/
    chmod + "${YELLOW}  Use: sudo micro $CONFIG_DIxNARY_PAi

# Parcho -ea servtente
if ss-actiuieE_NAME; thenciar${NC}"
    e
    YELLOW}Parando serviço existente...${NC}"nt.json antes de ini
    syNFIG_DIR/stemctl stop $SERVICE_NAMEconfig-clie
    sleep 1
fi $CO
LOW}  IMPORTANTE: Edite
# Copia binárioYEL
echo -e "${CYAN}Copiando binário para $INSTALL_DINCR..}"
    echo -e "${.${NC}"
cp "$BINARY_PATH" "$INSTALL_DIR/procs"ão criada${
ecpy-client"e "${GREEN}✓ Bináurarção padr${NC}
chmod +x "$INSTAL888}✓ Config
}
    echo -e "${YELLOWEOF
    
L_DIR/procspy-clientho -io instalado"
a dirórioport": 8s": "localhost",
    "api_
etdiretório.$
# Criecho -e "${CYAIR"N}Criando s..{NC}"
    "api_hostmkdir -p "$CONFIG_Ddor.com/procspy",

mkdir -tver_url": "https://seu-servierval": 5,
    "serp "R"
ec  "debug": false,
    "inho -"$LOG_DI$LOEEN}✓ DG",
  _DIR"
chmod 755 Re "${GRison" <<EOFretórios criados${NC
{ "$LOG_DI
    "user": "$(whoami)",
    "log_path":}"${NC}"
    onfig-client.j
    cat > "$CONFIG_DIR/c   # Cria configuração padrão
 
# Copia 
ou cria "$CONFIG_Pconfiguração}"ão padrão...
 quivo de configuraçãoiando configuraç não encontrado em: $CONFIG_PATH${NC}"
    echo -e "${CYAN}Cr
if [ -f ATH" ]; theno✓ Configuração copiada${NC}"
else -e "${YELLOW}Aviso: Ar
    echonfiguração: $CONFH${NCI
    cp "$CONFIG_PATH" "$CONFIG_DIR/config-client.json"
    echo -e "${GREEN}   echo -e "${CYAN}Copiando c deG_PAT