package main

import (
       "encoding/json"
        "fmt"

        "github.com/hyperledger/fabric/core/chaincode/shim"
        pb "github.com/hyperledger/fabric/protos/peer"
)

type Certificate struct {
        CertificateId           string  `json:"CertificateId"`
        FileHash                string  `json:"FileHash"`
        FilePath                string  `json:"FilePath"`
        DataHash                string  `json:"DataHash"`
        CertType                string  `json:"CertType"`
        Status                  string  `json:"Status"`
        Certifier1              string  `json:"Certifier1"`
        Certifier2              string  `json:"Certifier2"`
        Certifier3              string  `json:"Certifier3"`
        StudentAck              string  `json:"StudentAck"`
        AllValues               string  `json:"AllValues"`
        CurrentOwner            string  `json:"CurrentOwner"`
        TransferTo              string  `json:"TransferTo"`
}

type CertificateChaincode struct {
}

func (t *CertificateChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
        return shim.Success(nil)
}

func (t *CertificateChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
        fmt.Println("chaincode_custom Invoke")
        function, args := stub.GetFunctionAndParameters()

        if function == "query" {
                return t.query(stub, args)
        } else if function == "addActive" {
                return t.addActive(stub, args)
        } else if function == "studentAcknowledgement" {
                return t.studentAcknowledgement(stub, args)
        } else if function == "updateCertificate" {
                return t.updateCertificate(stub, args)
        } else if function == "deleteCertificate" {
                return t.deleteCertificate(stub, args)
        }

        return shim.Error("Invalid invoke function name")
}

func (t *CertificateChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var A string 
        var err error

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting certificate Id to query")
        }

        A = args[0]

        Avalbytes, err := stub.GetState(A)
        if err != nil {
                jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
                return shim.Error(jsonResp)
        }

        if Avalbytes == nil {
                jsonResp := "{\"Error\":\"Nil amount\"}"
                return shim.Error(jsonResp)
        }

        return shim.Success(Avalbytes)
}

func (t *CertificateChaincode) addActive(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        fmt.Println("Add Certificate.. ")
        var certificate Certificate
        var err error

        isExistAsBytes, err := stub.GetState(args[0])
        
        if err != nil { 
         return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err) 
                } 
                
        else if isExistAsBytes == nil { 
                return "", fmt.Errorf("Asset not found: %s", args[0]) 
        }

        if isExistAsBytes != nil {
          return shim.Error(err.Error())
        }

        certificate.CertificateId = args[0]
        certificate.FileHash      = args[1]
        certificate.FilePath      = args[2]
        certificate.Certifier1    = args[3]
        certificate.Certifier2    = args[4]
        certificate.Certifier3    = args[5]
        certificate.AllValues     = args[6]
        certificate.CurrentOwner  = args[7]

        certificateAsBytes, err := json.Marshal(certificate)
        err = stub.PutState(certificate.CertificateId, certificateAsBytes)
        if err != nil {
                return shim.Error(err.Error())
        }

        return shim.Success(nil)
}

func (t *CertificateChaincode) studentAcknowledgement(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        fmt.Println("Student Acknowledgement...")
        var certificate Certificate
        var err error

        certificate.CertificateId = args[0]
        certificate.StudentAck = args[1]

        certificateAsBytes, err := json.Marshal(certificate)
        err = stub.PutState(certificate.CertificateId, certificateAsBytes)

        if err != nil {
                return shim.Error(err.Error())
        }

        return shim.Success(nil)
}

func (t *CertificateChaincode) updateCertificate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        fmt.Println("Update Certificate")
        var err error
        certificateId := args[0]
        fileHash := args[1]
        filePath := args[2]
        certifier1 := args[3]
        certifier2 := args[4]
        certifier3 := args[5]
        allValues := args[6]

        certificateAsBytes, err := stub.GetState(certificateId)
        if err != nil {
                return shim.Error("Failed to get Certificate:" + err.Error())
        } else if certificateAsBytes == nil {
                return shim.Error("Certificate does not exist")
        }
        certificateToTransfer := Certificate{}
        err = json.Unmarshal(certificateAsBytes, &certificateToTransfer)
        if err != nil {
                return shim.Error(err.Error())
        }
		certificateToTransfer.FileHash = fileHash
		certificateToTransfer.FilePath = filePath
       	certificateToTransfer.Certifier1 = certifier1
       	certificateToTransfer.Certifier2 = certifier2
       	certificateToTransfer.Certifier3 = certifier3
       	certificateToTransfer.AllValues = allValues
        certificateJSONAsBytes, _ := json.Marshal(certificateToTransfer)
        err = stub.PutState(certificateId, certificateJSONAsBytes)
        if err != nil {
                return shim.Error(err.Error())
        }
        return shim.Success(nil)
}

func (t *CertificateChaincode) deleteCertificate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        if len(args) < 1 {
                return shim.Error("Incorrect number of arguments. Expectung 1")
        }

        Id := args[0]
        fmt.Printf("Deleting certificate %v", Id)
        valBytes, err := stub.GetState(Id)
        if err != nil {
                
                return shim.Error(err.Error())
        }

        if valBytes == nil {
        
                return shim.Error("Certificate '" + Id + "' not found.")
        }

        err = stub.DelState(Id)
        if err != nil {
			return shim.Error("Failed to delete state")
        }

         return shim.Success(nil)
}

func main() {
        err := shim.Start(new(CertificateChaincode))
        if err != nil {
                fmt.Printf("Error starting Certificate chaincode: %s", err)
        }
}
