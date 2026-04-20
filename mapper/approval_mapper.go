package mappers

import (
	"hrms_go/constant"
	"hrms_go/dto/approval"
	"hrms_go/models"
	"hrms_go/models/base"

	"github.com/google/uuid"
)

func ToApprovalTemplateDetail(dto approval.PostCreateApprovalTemplateDetailDto) (*models.ApprovalTemplateDetail, error) {
	headerId, err := uuid.Parse(dto.ApprovalTemplateHeaderId)
	if err != nil {
		return nil, err
	}

	approverBy, err := uuid.Parse(dto.ApproverBy)
	if err != nil {
		return nil, err
	}

	return &models.ApprovalTemplateDetail{
		ApprovalTemplateHeaderId: headerId,
		ApproverBy:               approverBy,
		SequenceNumber:           dto.SequenceNumber,
	}, nil
}

func ToApprovalDetail(dto approval.PostApproveDto) (models.ApprovalDetail,error)  {
	approvalHeaderId, err := uuid.Parse(dto.ApprovalHeaderId)
	if err != nil {
		return models.ApprovalDetail{}, err
	}

	approverBy, err := uuid.Parse(dto.ApproverBy)
	if err != nil {
		return models.ApprovalDetail{}, err
	}

	remark := ""
	if dto.Remark != nil{
		remark = *dto.Remark
	}

	return models.ApprovalDetail{
		ApprovalHeaderId: approvalHeaderId,
		ApprovalStatus: constant.ApprovalStatus(dto.ApprovalStatus),
		ApproverBy: approverBy,
		Remark: remark,
	},nil
}

func ToApprovalHeader(dto approval.CreateApprovalDto,approvalTemplateHeader models.ApprovalTemplateHeader) (models.ApprovalHeader,error)  {
	requesterBy, err := uuid.Parse(dto.RequesterBy)
	if err != nil {
		return models.ApprovalHeader{}, err
	}

	return models.ApprovalHeader{
		ApprovalTemplateHeaderId: approvalTemplateHeader.ApprovalTemplateHeaderId,
		ApprovalDocId: dto.ApprovalDocId,
		RequesterBy: requesterBy,
		AuditFields: base.AuditFields{
			CreatedBy: dto.RequesterBy,
		},
	},nil
}