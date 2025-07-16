package quotes

import (
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

// SOURCE: https://dev.freterapido.com.br/common/tipos_de_volumes/
var CategoryMap = map[int]string{
	1:   "Abrasivos",
	2:   "Adubos / Fertilizantes",
	3:   "Alimentos perecíveis",
	4:   "Artigos para Pesca",
	5:   "Auto Peças",
	6:   "Bebidas / Destilados",
	7:   "Brindes",
	8:   "Brinquedos",
	9:   "Calçados",
	10:  "CD / DVD / Blu-Ray",
	11:  "Combustíveis / Óleos",
	12:  "Confecção",
	13:  "Cosméticos",
	14:  "Couro",
	15:  "Derivados Petróleo",
	16:  "Descartáveis",
	17:  "Editorial",
	18:  "Eletrônicos",
	19:  "Eletrodomésticos",
	20:  "Embalagens",
	21:  "Explosivos / Pirotécnicos",
	22:  "Medicamentos",
	23:  "Ferragens",
	24:  "Ferramentas",
	25:  "Fibras Ópticas",
	26:  "Fonográfico",
	27:  "Fotográfico",
	28:  "Fraldas / Geriátricas",
	29:  "Higiene",
	30:  "Impressos",
	31:  "Informática / Computadores",
	32:  "Instrumento Musical",
	33:  "Livro(s)",
	34:  "Materiais Escolares",
	35:  "Materiais Esportivos",
	36:  "Materiais Frágeis",
	37:  "Material de Construção",
	38:  "Material de Irrigação",
	39:  "Material Elétrico / Lâmpada(s)",
	40:  "Material Gráfico",
	41:  "Material Hospitalar",
	42:  "Material Odontológico",
	43:  "Material Pet Shop",
	44:  "Material Veterinário",
	45:  "Móveis montados",
	46:  "Moto Peças",
	47:  "Mudas / Plantas",
	48:  "Papelaria / Documentos",
	49:  "Perfumaria",
	50:  "Material Plástico",
	51:  "Pneus e Borracharia",
	52:  "Produtos Cerâmicos",
	53:  "Produtos Químicos Não Classificados",
	54:  "Produtos Veterinários",
	55:  "Revistas",
	56:  "Sementes",
	57:  "Suprimentos Agrícolas / Rurais",
	58:  "Têxtil",
	59:  "Vacinas",
	60:  "Vestuário",
	61:  "Vidros / Frágil",
	62:  "Cargas refrigeradas/congeladas",
	63:  "Papelão",
	64:  "Móveis desmontados",
	65:  "Sofá",
	66:  "Colchão",
	67:  "Travesseiro",
	68:  "Móveis com peças de vidro",
	69:  "Acessórios de Airsoft / Paintball",
	70:  "Acessórios de Pesca",
	71:  "Simulacro de Arma / Airsoft",
	72:  "Arquearia",
	73:  "Acessórios de Arquearia",
	74:  "Alimentos não perecíveis",
	75:  "Caixa de embalagem",
	76:  "TV / Monitores",
	77:  "Linha Branca",
	78:  "Vitaminas / Suplementos nutricionais",
	79:  "Malas / Mochilas",
	80:  "Máquina / Equipamentos",
	81:  "Rações / Alimento para Animal",
	82:  "Artigos para Camping",
	83:  "Pilhas / Baterias",
	84:  "Estiletes / Materiais Cortantes",
	85:  "Produto Químico classificado",
	86:  "Limpeza",
	87:  "Extintores",
	88:  "Equipamentos de Segurança / API",
	89:  "Utilidades domésticas",
	90:  "Acessórios para celular",
	91:  "Toldos",
	92:  "Pisos (cerâmica) / Revestimentos",
	93:  "Artesanatos (sem vidro)",
	94:  "Quadros / Molduras",
	95:  "Porta / Janelas (sem vidro)",
	96:  "Placa de Energia Solar",
	97:  "Materiais hidráulicos / Encanamentos",
	98:  "Pia / Vasos",
	99:  "Bijuteria",
	100: "Joia",
	101: "Refrigeração Industrial",
	102: "Cocção Industrial",
	103: "Utensílios industriais",
	104: "Maquina de algodão doce",
	105: "Maquina de chocolate",
	106: "Estufa térmica",
	107: "Equipamentos de cozinha industrial",
	108: "Tapeçaria / Cortinas / Persianas",
	109: "Acessório para decoração (com vidro)",
	110: "Acessório para decoração (sem vidro)",
	111: "Acessórios automotivos",
	112: "Acessórios para bicicleta",
	113: "Artesanatos (com vidro)",
	114: "Bicicletas (desmontada)",
	115: "Cama / Mesa / Banho",
	116: "Chapas de madeira",
	117: "Manequins",
	118: "Portas / Janelas (com vidro)",
	119: "Torneiras",
	120: "Vasos de polietileno",
	121: "Chip de celular",
	122: "Celulares e Smartphones",
	123: "Telefonia Fixa e Sem Fio",
	124: "Portáteis industriais",
	125: "Eletrodomésticos industriais",
	126: "Expositor industrial",
	127: "Maçanetas",
	128: "Bebedouros e Purificadores",
	129: "Narguilés",
	130: "Acessórios para Narguilés",
	131: "Tabacaria",
	132: "Acessórios para Tabacaria",
	133: "Banheira Acrílico",
	134: "Banheira de Aço Esmaltada",
	135: "Banheira Fibra de Vidro",
	136: "Caixa Plástica",
	137: "Cartucho de Gás",
	138: "Equipamento oftalmológico",
	139: "Material oftalmológico",
	140: "Óculos e acessórios",
	141: "Material de laboratório",
	142: "Tinta",
	143: "Produtos de SexShop",
	144: "Lentes de contato",
	145: "Armações de óculos",
	146: "Caixa d'água (Completa)",
	147: "Tampa de Caixa d'água",
	148: "Selas e Arreios de montaria",
	149: "Acessórios de montaria",
	150: "Eletroportáteis",
	151: "Equipamentos para solda",
	152: "Artigos de festas",
	153: "Relógios",
	154: "Material de Jardinagem",
	155: "Acessório infantil",
	156: "Banheira infantil",
	157: "Cadeirinha para automóvel",
	158: "Carrinho de bebê",
	159: "Móveis infantis",
	999: "Outros",
}

type Quote struct {
	ID        uuid.UUID       `ch:"id"`
	Name      string          `ch:"name"`
	Service   string          `ch:"service"`
	Deadline  uint8           `ch:"deadline"`
	Price     decimal.Decimal `ch:"price"`
	CreatedAt time.Time       `ch:"timestamp"`
}

type RequestQuoteRecipientAddress struct {
	Zipcode string `json:"zipcode"`
}

type RequestQuoteRecipient struct {
	Address RequestQuoteRecipientAddress `json:"address"`
}

type RequestQuoteVolume struct {
	Category      int             `json:"category"`
	Amount        int             `json:"amount"`
	UnitaryWeight float64         `json:"unitary_weight"`
	Price         decimal.Decimal `json:"price"`
	Sku           string          `json:"sku"`
	Height        float64         `json:"height"`
	Width         float64         `json:"width"`
	Length        float64         `json:"length"`
}

type RequestQuote struct {
	Recipient RequestQuoteRecipient `json:"recipient"`

	Volumes []RequestQuoteVolume `json:"volumes"`
}

// ValidateCategories returns first invalid category index, -1 if all are valid
func (req RequestQuote) ValidateCategories() int {
	for idx, vol := range req.Volumes {
		if _, ok := CategoryMap[vol.Category]; !ok {
			return idx
		}
	}

	return -1
}

// ValidateDimensions returns first volume index, -1 if all are valid
func (req RequestQuote) ValidateDimensions() int {
	for idx, vol := range req.Volumes {
		if vol.Height <= 0 || vol.Width <= 0 || vol.Length <= 0 {
			return idx
		}
	}

	return -1
}

// ValidateAmount returns first volume index, -1 if all are valid
func (req RequestQuote) ValidateAmount() int {
	for idx, vol := range req.Volumes {
		if vol.Amount <= 0 {
			return idx
		}
	}

	return -1
}

// ValidatePrice returns first volume index, -1 if all are valid
func (req RequestQuote) ValidatePrice() int {
	for idx, vol := range req.Volumes {
		if vol.Price.IsZero() || vol.Price.IsNegative() {
			return idx
		}
	}

	return -1
}

// ValidateWeight returns first volume index, -1 if all are valid
func (req RequestQuote) ValidateWeight() int {
	for idx, vol := range req.Volumes {
		if vol.UnitaryWeight <= 0 {
			return idx
		}
	}

	return -1
}

// MustParseRecipientZipcode returns parsed recipient zipcode or -1 if invalid
func (req RequestQuote) MustParseRecipientZipcode() int64 {
	value, err := strconv.ParseInt(req.Recipient.Address.Zipcode, 10, 64)
	if err != nil {
		return -1
	}

	return value
}

// ParseRecipientZipcode returns parsed recipient zipcode or QuoteInvalidZipcode if invalid
func (req RequestQuote) ParseRecipientZipcode() (int64, error) {
	value := req.MustParseRecipientZipcode()
	if value == -1 {
		return -1, QuoteInvalidZipcode{
			Message: "Invalid recipient zipcode",
		}
	}

	return value, nil
}
