# Storm Surge Workflow

#set env variables

export FISSION_URL=http://$(kubectl --namespace fission get svc controller -o=jsonpath='{..ip}');
export FISSION_ROUTER=$(kubectl --namespace fission get svc router -o=jsonpath='{..ip}');

## set go path variables
# export GOPATH=$HOME/goworkspace;
# export PATH=$PATH:$HOME/goworkspace/bin;

# Build form-req-processor
# go-function-build ./form-req-processor/form-data-transformer.go;
fission function update --name transform-data --env go-env --deploy ./form-req-processor/function.so ;

# go-function-build ./go-wunderground-api/weather-api.go;
fission function update --name get-weather-data-v2 --env go-env --deploy ./go-wunderground-api/function.so;

# go-function-build ./weather-risk-api/compute-weather-risk.go;
fission function update --name get-weather-risk --env go-env --deploy ./weather-risk-api/function.so;

# go-function-build ./update-ticket-payload/update-ticket-data.go;
fission function update --name add-risk-data-to-ticket --env go-env --deploy ./update-ticket-payload/function.so; 

# go-function-build ./register-ticket/register-ticket.go;
fission function update --name register-ticket-v2 --env go-env --deploy  ./register-ticket/function.so; 

