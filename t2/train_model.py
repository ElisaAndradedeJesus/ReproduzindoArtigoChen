import numpy as np
import pandas as pd
from xgboost import XGBClassifier
from sklearn.model_selection import train_test_split
# Biblioteca para exportação em formato PMML compatível com sistemas nativos
from nyoka import xgboost_to_pmml 

# 1. Simulação das colunas baseadas na Tabela I do artigo (Excluindo IPs das features)
features = [
    "Timestamp", "Source_Port", "Destination_Port", "Protocol", 
    "Flow_Packets_s", "Flow_Bytes_s", "ACK_Flag_Count", "SYN_Flag_Count", 
    "RST_Flag_Count", "URG_Flag_Count", "CWR_Flag_Count", "Packet_Length_Mean"
]

# Criando dados sintéticos para demonstração do pipeline
np.random.seed(42)
X_dummy = np.random.rand(1000, len(features))
y_dummy = np.random.choice([0, 1], size=1000, p=[0.8, 0.2]) # 0: Normal, 1: DDoS

df = pd.DataFrame(X_dummy, columns=features)

X_train, X_test, y_train, y_test = train_test_split(df, y_dummy, test_size=0.2)

# 2. Inicialização e treinamento do algoritmo XGBoost escolhido no artigo
model = XGBClassifier(
    max_depth=6, 
    learning_rate=0.1, 
    n_estimators=100, 
    eval_metric="logloss"
)
model.fit(X_train, y_train)

# Avaliação simplificada do pipeline
accuracy = model.score(X_test, y_test)
print(f"Modelo treinado com sucesso! Acurácia simulada: {accuracy * 100:.2f}%")

# 3. Exportação em formato de marcação padrão (PMML) para o eBPF ler no User Space
# Isso cria a estrutura lida dinamicamente no espaço de usuário do Go/C++
xgboost_to_pmml(model, features, "target", "xgboost_ddos_model.pmml")
print("Arquivo 'xgboost_ddos_model.pmml' gerado e pronto para implantação.")
