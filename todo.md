// Domain model

job
  (id)
  task_form

  verification_form
  verification_module
  verification_agreement_count
  verification_score_threshold
  verification_prompt
  funding_verification_reward
  verification_limit

  task
    (job_id, id)
    data
    is_active
    assignment_count
    pending_count
    accepted_count
    verification_count

//- task flow
assignment
  (job_id, task_id, id)
  user_id

responses
  (job_id, task_id, id)
  worker_id
  value
  is_accepted
  verifications_count

//- verification flow
verification_assignment
  (job_id, task_id, response_id, id)
  user_id  

verification_responses
  (job_id, task_id, response_id, id)
  worker_id
  verifier_id
  value
  reason

verificaiton_eligibility
  job_id
  user_id


// modules --------------------------------
verification-settings
  get(jobId)
  create(settings)
  
  // structures
  verification-settings
    verficaiton_type
      automatic
      manual
    
    manual
      form
    
    verification_agreement_count
    verification_prompt
    funding_verification_reward
    verification_limit


verification-assignment
  assign(jobId, userId)
  unassign(jobId, userId)

verification-verify
  verify(userId, responseId, result, reason)
  



func Verify(userID uint64, responseID uint64, result bool, reason string) (interface{}, error) {

	// Ensure verification is assigned
	isAssigned, err := verificationassignmentsvc.CheckAssignedVerification(userID, responseID)

	response, err := responsesvc.GetResponse(responseID)

	settings, err := settingssvc.GetSettings(response.JobID)

	// Get job's verification module
	verificationModule := verification.VerificationModules[settings.VerificationModule]
	
	vresult, err = verificationModule.Verify(
		userID, 
		response, 
		result,
		reason
	)

	if err != nil {
		return nil, err
	}
	return result, nil
}

