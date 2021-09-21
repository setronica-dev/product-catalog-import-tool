#!/bin/bash
function create_config {
  DIR=$1
  cat >./service.yaml <<EOF
product:
  source: ${DIR}/data/source/products/
  report: ${DIR}/data/result/report/
  sent: ${DIR}/data/source/processed/products/
  in_progress: ${DIR}/data/source/inprogress/
  success_result: ${DIR}/data/result/sent/
  mapping: ${DIR}/data/mapping/mapping.yaml
  ontology: ${DIR}/data/ontology/rules.csv

offer:
  source: ${DIR}/data/source/offers/
  sent: ${DIR}/data/source/processed/offers/

offer_item:
  source: ${DIR}/data/source/offerItems/
  success_result: ${DIR}/data/result/sent/offerItems/
  report: ${DIR}/data/result/report/
  sent: ${DIR}/data/source/processed/offerItems/

xlsx_settings:
  source: ${DIR}/data/source/
  sent: ${DIR}/data/source/processed/
  sheet:
    products:
      name: "Products"
      header_rows_to_skip: 2
    offers:
      name: "Offers"
      header_rows_to_skip: 2
    attributes:
      name: "Attributes"
      header_rows_to_skip: 2
    offer_items:
      name: "Prices"
      header_rows_to_skip: 2

tradeshift_api:
  # set Tradeshift API parameters from API Access To Own Account in Tradeshift pannel
  base_url:
  consumer_key:
  consumer_secret:
  token:
  token_secret:
  tenant_id:
  currency:
  file_locale:
  recipients:
    - id:
      name:
EOF
}

usage() { echo "Usage: $0 [-d <targed dir>]" 1>&2; exit 1; }

while getopts ":d:" o; do
    case "${o}" in
        d)
            DIR=${OPTARG}
            ;;
        *)
            usage
            ;;
    esac
done
shift $((OPTIND-1))
if [ -z "${DIR}" ]; then
    usage
fi

go get ./...

mkdir -p $DIR
go build -o ./product-catalog-import-tool

cp ./product-catalog-import-tool $DIR/product-catalog-import-tool
create_config $DIR
cp ./service.yaml $DIR/
cp -R ./data $DIR/
